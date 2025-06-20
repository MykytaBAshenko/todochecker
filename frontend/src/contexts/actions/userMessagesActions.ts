import { createAsyncThunk } from '@reduxjs/toolkit';
import * as apiClient from '../../api-client';
import { showToast } from '../slices/toastSlice';

export const getAllUserMessages = createAsyncThunk(
  'usermessages/all',
  async (_, { rejectWithValue }) => {
    try {
      const data = await apiClient.getAllUserMessages();
      return data;
    } catch (error: any) {
      rejectWithValue(error)
    }
  }
);

