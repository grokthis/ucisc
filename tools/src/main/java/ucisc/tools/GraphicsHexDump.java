package ucisc.tools;

import java.io.*;

public class GraphicsHexDump {
  private int size;
  private String mode;
  private int xResolution;
  private int yResolution;
  private int charsetAddress;
  private File imageFile;
  private String outputFile;
  private int[] words;
  private int colorWidth;

  public GraphicsHexDump(
      int size,
      String mode,
      String resolution,
      int charsetAddress,
      String outFile,
      String imageFile
  ) {
    this.size = size;
    this.mode = mode;
    String[] dimensions = resolution.split("x");
    this.xResolution = Integer.parseInt(dimensions[0]);
    this.yResolution = Integer.parseInt(dimensions[1]);
    this.charsetAddress = charsetAddress;
    this.imageFile = new File(imageFile);
    this.outputFile = outFile;
    words = new int[size];
    if ("text".equals(mode)) {
      if (xResolution * yResolution * 3 < size) {
        colorWidth = 16;
      } else if (xResolution * yResolution * 2 < size) {
        colorWidth = 8;
      } else {
        throw new RuntimeException(
            "Text resolution ("
                + dimensions
                + ") too larger for memory ("
                + size
                + ")");
      }
    } else {
      throw new RuntimeException("Mode not implemented: " + mode);
    }
  }

  public void write() throws IOException {
    for (int j = 0; j < 6; j++) {
      String fileName = String.format(outputFile, j);
      BufferedWriter writer = new BufferedWriter(new FileWriter(fileName));
      for (int i = j * 16384; i < (j + 1) * 16384; i++) {
        String hex = String.format("%04X ", words[i]);
        if (i > 0 && i % 16 == 0) {
          writer.newLine();
        }
        writer.write(hex);
      }
      writer.close();
    }
  }

  public void generateData() throws Exception {
    if ("text".equals(mode)) {
      clearText();
      drawText();
      drawTextForeground();
      drawTextBackground();
      drawCharset();
    }
  }

  private void clearText() {
    // Set all visible characters to "space"
    for (int i = 0; i < xResolution * yResolution; i++) {
      words[i] = 0x20;
    }
  }

  private void drawText() throws IOException {
    if (this.imageFile == null) {
      return;
    }
    BufferedReader reader;
    if (imageFile.exists()) {
      reader = new BufferedReader(new FileReader(imageFile));
    } else {
      System.err.println("Image file is missing: " + imageFile);
      throw new RuntimeException("Image file is missing: " + imageFile);
    }

    String line;
    int y = 0;
    while (y < this.yResolution && (line = reader.readLine()) != null) {
      for (int x = 0; x < xResolution && x < line.length(); x++) {
        words[y * xResolution + x] = line.codePointAt(x);
      }
      y += 1;
    }
  }

  private int randomColor(int charPos, int totalChars, boolean fg) {
    int mask = 0xFFFF;
    int shift = 0;
    if (fg) {
//      while (totalChars < 0x8000) {
//        totalChars = totalChars << 1;
//        mask = mask << 1;
//        shift = shift + 1;
//      }
//      mask = mask ^ 0xFFFF;
//      return (charPos << shift) | (mask & (int)(Math.random() * 0xFFFF));
      return ((int)(Math.random() * 0x10000)) & 0xFFFF;
    } else {
      int channel = (int)(((charPos / xResolution) / (double)yResolution) * 0x1F);
      return (channel << 11) | (((channel << 1) | 0x1) << 5) | channel;//(channel << 10) | (channel << 5) | channel;
    }
  }

  private void drawTextForeground() {
    int foregroundStart;
    int totalChars = xResolution * yResolution;
    if (totalChars < 16384) {
      // 16-bit color depth
      foregroundStart = 16384;
    } else if (totalChars < 24576) {
      // 16-bit color depth
      foregroundStart = 24576;
    } else {
      // 8-bit color depth, fg and bg are in the same word
      foregroundStart = 40960;
    }

    for (int i = 0; i < totalChars; i++) {
      int color = randomColor(i, totalChars, true);
      if (totalChars < 16384) {
        //color = 0x07E0; // blue
      } else if (totalChars < 24576) {
        //color = 0x07E0; // blue
      } else {
        // 8-bit color depth, fg and bg are in the same word
        //color = 0x06AA; // dark blue fg, yellow bg
        color = (0xFF00 & randomColor(i, totalChars, true)) | (0xFF & randomColor(i, totalChars, true));
      }
      words[foregroundStart + i] = color;
    }
  }

  private void drawTextBackground() {
    //int color = 0xFFE0;
    int backgroundStart;
    int totalChars = xResolution * yResolution;
    if (totalChars < 16384) {
      backgroundStart = 32768;
    } else if (totalChars < 24576) {
      backgroundStart = 49152;
    } else {
      return;
    }
    // Set background color to uniform hard coded value for now
    for (int i = 0; i < totalChars; i++) {
      words[backgroundStart + i] = randomColor(i, totalChars, false);
    }
  }

  private void drawCharset() {
    Charset charset = new Charset();
    for (int i = charsetAddress; i < size; i++) {
      words[i] = charset.getWord(i - charsetAddress);
    }
  }
}
