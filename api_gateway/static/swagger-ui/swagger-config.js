window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: [
      {url :"/swagger/product_service.json", name: "Product Service"},
      {url :"/swagger/customer_service.json", name: "Customer Service"},
      {url :"/swagger/cart_service.json", name: "Cart Service"},
      {url :"/swagger/order_service.json", name: "Order Service"},
      {url :"/swagger/payment_service.json", name: "Payment Service"},
      {url :"/swagger/shipment_service.json", name: "Shipment Service"}

    ],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};