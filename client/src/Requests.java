import java.io.IOException;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

public class Requests {

public String run(String url) throws IOException {
  OkHttpClient client = new OkHttpClient();
  Request request = new Request.Builder().url(url).build();

  try (Response response = client.newCall(request).execute()) {
    return response.body().string();
  }
}
}