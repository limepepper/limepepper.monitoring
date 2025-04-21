#!/bin/bash

# Find USB mic by vendor/product ID
USB_DEV_PATH=$(find /sys/bus/usb/devices -name idVendor -print0 | xargs -0 -I{} sh -c 'grep -l 0d8c "{}" && echo -n "{}" | xargs -0 dirname' | xargs -0 -I{} find {} -name idProduct -print0 | xargs -0 -I{} sh -c 'grep -l 013c "{}" && echo -n "{}" | xargs -0 dirname')

if [ -z "$USB_MIC_CARD" ]; then
  echo "Error: USB PnP Sound Device not found" >&2
  exit 1
fi

echo "Using sound card: $USB_MIC_CARD" >&2

# Find the sound card number associated with this USB device
CARD_NUM=$(find $USB_DEV_PATH -path "*/sound/card*" -name "cardnum" -exec cat {} \; 2>/dev/null ||
          find /sys/class/sound/card* -name "device" -exec readlink -f {} \; | grep -n "$USB_DEV_PATH" | cut -d: -f1)

if [ -z "$CARD_NUM" ]; then
  # Alternative: try to find by name if path approach fails
  CARD_NUM=$(arecord -l | grep "USB PnP Sound Device" | head -1 | sed -r 's/card ([0-9]+).*/\1/')
  if [ -z "$CARD_NUM" ]; then
    echo "Error: Could not find sound card number" >&2
    exit 1
  fi
fi

echo "Using sound card: $CARD_NUM" >&2

# Capture sox stats output
stats_output=$(sox -t alsa hw:$CARD_NUM,0 -n trim 0 3 stats 2>&1)

# Rest of script as before...
