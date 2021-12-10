package ucisc.tools;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * General utilities command line interface for ucisc tools.
 */
public class Utilities {
  public static void printUsage() {
    List<String> usage = new ArrayList<>();
    usage.add("Usage: ucisc-tools <command> [command args...]");
    usage.add("Commands:");
    usage.add("  g-hex: Generates a graphics memory hex dump for initializing video memory");
    usage.add("    args: <size> <mode> <resolution> <charset-address> <out-file> [<image-file>]");
    usage.add("      size: total video memory size in bytes (e.g. 81920 for 80K)");
    usage.add("      mode: graphics mode - supports 'text' or 'bitmap'");
    usage.add("      resolution: pixel or character resolution (e.g. 1280x960 or 80x60)");
    usage.add("      charset-address: starting character set address (e.g. 77824)");
    usage.add("      out-file: file to output hex dump to");
    usage.add("      image-file: image to text file with initial memory data");

    usage.forEach(line -> {
      System.out.println(line);
    });
  }

  public static void main(String[] args) throws Exception {
    if (args.length < 1) {
      printUsage();
      System.exit(1);
    } else if (args[0].equals("g-hex")) {
      parseGraphicsHex(Arrays.copyOfRange(args, 1, args.length));
      return;
    }
    System.exit(1);
  }

  public static void parseGraphicsHex(String[] args) throws Exception {
    if (args.length < 5) {
      printUsage();
      System.exit(1);
    }
    int size = Integer.parseInt(args[0]);
    String mode = args[1];
    String resolution = args[2];
    int charsetAddress = Integer.parseInt(args[3]);
    String outFile = args[4];
    String imageFile = null;
    if (args.length > 5) {
      imageFile = args[5];
    }

    GraphicsHexDump dumper = new GraphicsHexDump(
        size,
        mode,
        resolution,
        charsetAddress,
        outFile,
        imageFile
    );
    dumper.generateData();
    dumper.write();
  }
}
