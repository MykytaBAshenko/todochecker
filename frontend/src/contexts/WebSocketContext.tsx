import React, { createContext, useContext, useEffect, useRef, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { RootState } from './store';
import { showToast } from './slices/toastSlice';
// import { addMessage } from './slices/messagesSlice';

import { addGroupedMessage, deleteConversationGroupedMessage } from './slices/userMessagesSlice';
import { addNewUserMessanger, deleteUserMessanger } from './slices/userMessagesUsersSlice';

type WebSocketContextType = WebSocket | null;

const WebSocketContext = createContext<WebSocketContextType>(null);
const serverWS = import.meta.env.VITE_API_BASE_SOCKET_PROTOCOL + import.meta.env.VITE_API_BASE_URL || "";

export const useWebSocket = (): WebSocketContextType => {
  return useContext(WebSocketContext);
};

interface Props {
  children: React.ReactNode;
}

export const WebSocketProvider: React.FC<Props> = ({ children }) => {
  const socketRef = useRef<WebSocket | null>(null);
  const [socketError, setSocketError] = useState<string | null>(null); // Track socket errors
  const user = useSelector((state: RootState) => state.auth.user);
  const dispatch = useDispatch();

  const connectWebSocket = () => {
    const storedToken = localStorage.getItem('token');
    
    if (!storedToken) {
      console.warn('No token found in localStorage');
      return;
    }

    const ws = new WebSocket(`${serverWS}/ws?token=${storedToken}`);
    socketRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket connected');
      setSocketError(null); // Reset any previous error
    };

    ws.onmessage = (event) => {
      try {
        const {data ,type} = JSON.parse(event.data);
        console.log(data)
    
        if (type === "incoming_message") {
          dispatch(addGroupedMessage({ message: data , currentUserId: user.id })); // assuming `message.payload` is the message object
          dispatch(showToast({ type: "SUCCESS", message: `New message from ${data.message.MessageSender}` }));
        }

        if (type === "initial_message") {
          dispatch(addNewUserMessanger({message: data })); 
          dispatch(addGroupedMessage({message:data, currentUserId: user.id })); 
          // dispatch(showToast({ type: "SUCCESS", message: `New message from ${data.message.MessageSender}` }));
        }

        if (type === "delete_conversation") {
          dispatch(deleteUserMessanger({message: data })); 
          dispatch(deleteConversationGroupedMessage({message:data})); 
        }
      } catch (e) {
        console.log(e)
        console.error("Invalid JSON from server:", e);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      if (error instanceof ErrorEvent) {
        console.error('Error details:', error.message);
      }
      setSocketError('WebSocket error occurred. Retrying...'); // Update error state
    };

    ws.onclose = (event) => {
      console.log('WebSocket closed', event.code, event.reason);
      if (event.code !== 1000) {
        console.warn(`WebSocket closed with abnormal code: ${event.code}`);
        setSocketError(`WebSocket closed with code ${event.code}. Reconnecting...`); // Update error state
        setTimeout(() => {
          connectWebSocket(); // Try reconnecting after 3 seconds
        }, 3000);
      }
    };

    // Set up ping to keep the connection alive
    // const pingInterval = setInterval(() => {
    //   if (ws.readyState === WebSocket.OPEN) {
    //     ws.send('ping');
    //   } else {
    //     console.error('WebSocket is not open, cannot send ping');
    //   }
    // }, 30000); // Ping every 30 seconds

    return () => {
      clearInterval(pingInterval);
      if (socketRef.current) {
        socketRef.current.close();
        console.log('WebSocket connection closed during cleanup');
      }
    };
  };

  useEffect(() => {
    connectWebSocket();

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
        console.log('WebSocket connection closed during cleanup');
      }
    };
  }, [user]);

  return (
    <WebSocketContext.Provider value={socketRef.current}>
      {children}
      {socketError && <div style={{ color: 'red' }}>{socketError}</div>} {/* Display error message */}
    </WebSocketContext.Provider>
  );
};