public class Client {
  public static void main(String[] args) throws Exception {
    System.out.println("Making requests...");
    CRequest request = new CRequest();
    request.run("/api/auctions");
  }
}