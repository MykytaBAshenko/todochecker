import { RegisterFormData } from "./pages/Register";
import { SignInFormData } from "./pages/SignIn";
import {
  PaymentIntentResponse,
} from "../../backend/src/shared/types";
import axios from 'axios';
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "";
const axiosInstance = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

axiosInstance.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});


export const register = async (formData: RegisterFormData) => {
  const response = await axiosInstance.post('/auth/signup', formData);
  const { user, token } = response.data;
  localStorage.setItem('token', token);
  return user;
};

export const signIn = async (formData: SignInFormData) => {
    const response = await axiosInstance.post('/auth/login', formData);
    const { user, token } = response.data;
    localStorage.setItem('token', token);
    return user;
};

export const validateToken = async () => {
  const response = await axiosInstance.get(`${API_BASE_URL}/auth/validate-token`);
  const { user } = response.data;
  return user;
};

export const signOut = async () => {
  localStorage.removeItem('token');
};

export const getAllUserMessages = async () => {
  const response = await axiosInstance.get(`${API_BASE_URL}/usermessages/all`);
  console.log(response.data);
  return response.data;
};


export const getAllUserMessagesUsers = async () => {
  const response = await axiosInstance.get(`${API_BASE_URL}/usermessages/allusers`);
  console.log(response.data);
  return response.data;
};













export const createPaymentIntent = async (
  hotelId: string,
  numberOfNights: string
): Promise<PaymentIntentResponse> => {
  const response = await fetch(
    `${API_BASE_URL}/api/hotels/${hotelId}/bookings/payment-intent`,
    {
      credentials: "include",
      method: "POST",
      body: JSON.stringify({ numberOfNights }),
      headers: {
        "Content-Type": "application/json",
      },
    }
  );

  if (!response.ok) {
    throw new Error("Error fetching payment intent");
  }

  return response.json();
};

