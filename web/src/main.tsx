import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { Toaster } from "@/components/ui/sonner.tsx";
import { ThemeProvider } from "@/components/theme-provider.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ThemeProvider storageKey="vidlp-theme">
      <App />
      <Toaster richColors position="top-center" />
    </ThemeProvider>
  </React.StrictMode>,
);
