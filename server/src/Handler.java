import java.io.IOException;
import java.io.OutputStream;
import java.net.URI;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

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
    String response = "";
    switch(uri) {
      case "/api/auctions":
        System.out.println("Auctions");
        response = "auctions";
        break;
    }
    return response;
  }
}