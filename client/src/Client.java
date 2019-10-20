import java.io.IOException;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

public class Client {
  public static void main(String[] args) throws Exception {
    System.out.println("Making requests...");

    // Test URLs
    String[] testURLs = {
      "http://localhost:8080/api/auctions",
    };

    // Run tests
    for (int i = 0; i < testURLs.length; i++) {
      String res = sendRequest(testURLs[i]);
      System.out.println(res);
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