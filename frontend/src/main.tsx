// punto de entrada de React; monta el router.
import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import App from "./App";
import AuthPage from "./pages/AuthPage";
import VaultPage from "./pages/VaultPage";

const router = createBrowserRouter([
  { path: "/", element: <App />, children: [
    { index: true, element: <AuthPage /> },
    { path: "/vault", element: <VaultPage /> },
  ]},
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);
