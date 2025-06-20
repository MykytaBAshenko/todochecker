import React, { useContext, useState, useEffect } from "react";
import Toast from "../components/Toast";
import { useDispatch, useSelector } from "react-redux";
import { validateUserToken } from "./actions/authActions";
import { getAllUserMessages } from "./actions/userMessagesActions";
import { getAllUserMessagesUsers } from "./actions/userMessagesUsersActions";

import { RootState } from "./store";
import { hideToast } from "./slices/toastSlice";
import { loadStripe, Stripe } from "@stripe/stripe-js";
import { showToast } from "./slices/toastSlice";
const STRIPE_PUB_KEY = import.meta.env.VITE_STRIPE_PUB_KEY || "";
import { WebSocketProvider } from './WebSocketContext';
type ToastMessage = {
  message: string;
  type: "SUCCESS" | "ERROR";
};

type AppContext = {
  showToast: (toastMessage: ToastMessage) => void;
  user: any;
  signOut: Function;
};

// Create a context for backward compatibility
const AppContext = React.createContext<AppContext | undefined>(undefined);

export const AppContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const dispatch = useDispatch();
  const toast = useSelector((state: RootState) => state.toast);
  const user = useSelector((state: RootState) => state.auth.user);
  
  useEffect(() => {
    dispatch(validateUserToken());
    dispatch(getAllUserMessages());
    dispatch(getAllUserMessagesUsers());
    
  }, [!user]);


  return (
    <WebSocketProvider>
      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => dispatch(hideToast())}
        />
      )}
      {children}
    </WebSocketProvider>
  );
};
