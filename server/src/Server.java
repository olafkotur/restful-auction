import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpServer;

public class Server {
  public static void main(String[] args) throws Exception {
    int port = 8080;

    HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
    server.createContext("/", new Handler());
    server.start();

    System.out.println("Listening on Port " + port + "...");
  }
}