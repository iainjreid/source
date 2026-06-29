
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./app/app.js";
import "./styles/index.css";
import { GetRepos } from "./app/utils/api.js";

GetRepos().promise.then(({ repos }) => {
  createRoot(document.getElementById("root")!).render(
    <StrictMode>
      <App repos={repos} />
    </StrictMode>
  );
})
