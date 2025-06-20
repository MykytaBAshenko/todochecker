import { configureStore } from '@reduxjs/toolkit';
import authReducer from './slices/authSlice';
import toastReducer from './slices/toastSlice';
import userMessagesReducer from './slices/userMessagesSlice';
import userMessagesUsersReducer from './slices/userMessagesUsersSlice';



export const store = configureStore({
  reducer: {
    auth: authReducer,
    usermessages: userMessagesReducer,
    usermessagesusers: userMessagesUsersReducer,
    toast: toastReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;