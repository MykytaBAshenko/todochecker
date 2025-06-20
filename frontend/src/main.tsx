import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.scss";
import { Provider } from 'react-redux';
import { store } from './contexts/store.ts';
import { AppContextProvider } from "./contexts/AppContext.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider store={store}> 
      <AppContextProvider>
          <App />
      </AppContextProvider>
    </Provider>
  </React.StrictMode>
);
