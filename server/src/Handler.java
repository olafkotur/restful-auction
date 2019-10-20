import java.io.IOException;
import java.io.OutputStream;
import java.net.URI;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

import org.json.JSONObject;

class Handler implements HttpHandler {
  public void handle(HttpExchange ex) throws IOException {

    String requestMethod = ex.getRequestMethod();
    URI requestUri = ex.getRequestURI();

    System.out.println("\nMethod: " + requestMethod);
    System.out.println("URI: " + requestUri);

    String response = getResponse(requestUri.toString());
    ex.sendResponseHeaders(200, response.length());

    OutputStream os = ex.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }

  public String getResponse(String uri) {
    JSONObject res = new JSONObject();
    if (uri.contains("/api/auctions")) {
      res.put("code", 1);
      res.put("status", "success");
      res.put("message", "It's working! :)");
    }
    return res.toString();
  }
}