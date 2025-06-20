import { createAsyncThunk } from '@reduxjs/toolkit';
import * as apiClient from '../../api-client';

export const getAllUserMessagesUsers = createAsyncThunk(
    'usermessages/allusers',
    async (_, { rejectWithValue }) => {
      try {
        const data = await apiClient.getAllUserMessagesUsers();
        return data;
      } catch (error: any) {
        rejectWithValue(error)
      }
    }
  );