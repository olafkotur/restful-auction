public class Client {
  public static void main(String[] args) throws Exception {
    System.out.println("Making requests...");
    Requests request = new Requests();

    // Test URLs
    String[] testURLs = {
      "http://localhost:8080/api/auctions",
    };

    // Run tests
    for (int i = 0; i < testURLs.length; i++) {
      String res = request.run(testURLs[i]);
      System.out.println(res);
    }
  }
}