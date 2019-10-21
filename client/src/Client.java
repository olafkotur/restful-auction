import java.io.IOException;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

public class Client {
  public static void main(String[] args) throws IOException {
    // Test URLs
    String[] testURLs = {
      "http://localhost:8080/api/auctions",
    };

    runTests(testURLs);
  }

  public static void runTests(String[] testURLs) throws IOException {
    System.out.println("Attempting to run auction tests...");

    for (int i = 0; i < testURLs.length; i++) {
      String res = sendRequest(testURLs[i]);
      System.out.println("\nRequested: " + testURLs[i]);
      System.out.println("Response: " + res);
    }
  }

  public static String sendRequest(String url) throws IOException {
  OkHttpClient client = new OkHttpClient();
  Request request = new Request.Builder().url(url).build();

  try (Response response = client.newCall(request).execute()) {
    return response.body().string();
  }
}
}