import { createAsyncThunk } from '@reduxjs/toolkit';
import * as apiClient from '../../api-client';
import { SignInFormData } from '../../pages/SignIn';
import { RegisterFormData } from '../../pages/Register';
import { showToast } from '../slices/toastSlice';

type Error = {
    message: string | null;
};
  

export const loginUser = createAsyncThunk(
  'auth/login',
  async (credentials: SignInFormData, { dispatch, rejectWithValue }) => {
    try {
      const user = await apiClient.signIn(credentials);
      dispatch(showToast({ message: "Sign in Successful!", type: "SUCCESS" }));
      return user;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to login';
      dispatch(showToast({ message: errorMessage, type: "ERROR" }));
      return rejectWithValue(errorMessage);
    }
  }
);

export const registerUser = createAsyncThunk(
  'auth/register',
  async (userData: RegisterFormData, { dispatch, rejectWithValue }) => {
    try {
      const formData = { ...userData };
      formData.avatar = Array.from({ length: 8 }, () => Math.floor(Math.random() * 10)).join('');
      
      const user = await apiClient.register(formData);
      dispatch(showToast({ message: "Registration Success!", type: "SUCCESS" }));
      return user;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to register';
      dispatch(showToast({ message: errorMessage, type: "ERROR" }));
      return rejectWithValue(errorMessage);
    }
  }
);

export const validateUserToken = createAsyncThunk(
  'auth/validateToken',
  async (_, { dispatch, rejectWithValue }) => {
    try {
        const token = localStorage.getItem("token");
        console.log(token)
        if (token) {
            const user = await apiClient.validateToken();
            return user;
        }
        return null;
    } catch (error) {
        if (error?.message?.includes("401")) {
            localStorage.removeItem('token');
        }
        dispatch(showToast({ message: `Token validation failed: ${error.message}`, type: "ERROR" }));
        return rejectWithValue('Token validation failed');
    }
  }
);

export const logoutUser = createAsyncThunk(
  'auth/logout',
  async (_, { dispatch }) => {
    await apiClient.signOut();
    dispatch(showToast({ message: "Signed Out Successfully", type: "SUCCESS" }));
    return null;
  }
);