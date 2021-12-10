package ucisc.tools;

import org.junit.jupiter.api.Test;

/**
 * Tests for utilities command line interface.
 */
public class UtilitiesTest {
  @Test
  public void testRun() throws Exception {
    System.out.println("Hello from test.");
    Utilities.main(new String[0]);
  }
}
