import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { getAllUserMessagesUsers } from "../actions/userMessagesUsersActions";
import { UserMessage } from "./userMessagesUsersSlice";

export interface User {
  ID: number;
  Nickname: string;
  Email: string;
  Avatar?: string;

  // Add other user fields if needed
}

interface UsersState {
  groupedUsers: User[];
  isLoading: boolean;
  error: string | null;
}

const initialState: UsersState = {
  groupedUsers: [],
  isLoading: false,
  error: null,
};

const userMessagesUsersSlice = createSlice({
  name: "userMessagesUsers",
  initialState,
  reducers: {
    addNewUserMessanger(
      state,
      action: PayloadAction<{ message: UserMessage }>
    ) {
      const { message } = action.payload;
      console.log(message);
      return {
        ...state,
        groupedUsers: [...state.groupedUsers, message.Sender],
      };
    },
    deleteUserMessanger(
      state,
      action: PayloadAction<{ message: { userID: number | string } }>
    ) {
      const { userID } = action.payload.message;

      return {
        ...state,
        groupedUsers: state.groupedUsers.filter(user => user.ID !== userID),
      };
    },
  },
  extraReducers: (builder) => {
    builder.addCase(getAllUserMessagesUsers.pending, (state) => {
      state.isLoading = true;
      state.error = null;
    });
    builder.addCase(
      getAllUserMessagesUsers.fulfilled,
      (state, action: PayloadAction<User[]>) => {
        state.groupedUsers = action.payload;
        state.isLoading = false;
      }
    );
    builder.addCase(getAllUserMessagesUsers.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload as string;
    });
  },
});

export default userMessagesUsersSlice.reducer;
export const { addNewUserMessanger, deleteUserMessanger } = userMessagesUsersSlice.actions;
