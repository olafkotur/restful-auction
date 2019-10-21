import java.net.InetSocketAddress;
import com.sun.net.httpserver.HttpServer;
import java.util.Map;
import java.util.HashMap;
import java.io.IOException;
import java.io.OutputStream;
import java.net.URI;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.Headers;
import org.json.JSONObject;

public class Server {

  static final int PORT = 8080;

  public static void main(String[] args) throws IOException {
    HttpServer server = HttpServer.create(new InetSocketAddress(PORT), 0);
    server.createContext("/api", Server::handleRequest);
    server.start();

    System.out.println("Listening on Port " + PORT + "...");
  }

  public static void handleRequest(HttpExchange ex) throws IOException {
    printRequest(ex);

    String uri = ex.getRequestURI().toString();
    String response = getResponse(uri);

    ex.getResponseHeaders().set("Content-Type", "application/json;");
    ex.sendResponseHeaders(200, response.length());
    OutputStream os = ex.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }

  public static void printRequest(HttpExchange ex) {
    Headers headers = ex.getRequestHeaders();
    String method = ex.getRequestMethod().toString();
    String uri = ex.getRequestURI().toString();

    System.out.println("\n\n*** Headers ***");
    headers.entrySet().forEach(System.out::println);
    System.out.println("\nMethod: " + method);
    System.out.println("URI: " + uri);
  }

  public static String getResponse(String uri) {
    JSONObject res = new JSONObject();

    if (uri.equals("/api/auctions")) {
      res.put("code", 1);
      res.put("status", "success");
      res.put("message", "It's working! :)");
      return res.toString();
    }

    return "";
  }
}
