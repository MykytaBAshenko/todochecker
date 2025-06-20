import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { getAllUserMessages } from '../actions/userMessagesActions';
import { RootState } from '../store';
import { User } from './userMessagesUsersSlice'


export interface UserMessage {
  ID: number;
  MessageSender: number;
  MessageReceiver: number;
  MessageBody: string;
  CreatedAt: string;
  UpdatedAt?: string;
  Sender?: User;
  Reciver?: User;
}

interface GroupedMessagesState {
  groupedMessages: Record<number, UserMessage[]>;
  isLoading: boolean;
  error: string | null;
}

const initialState: GroupedMessagesState = {
  groupedMessages: {},
  isLoading: false,
  error: null
};

const userMessagesSlice = createSlice({
  name: 'userMessages',
  initialState,
  reducers: {
    addGroupedMessage(
      state,
      action: PayloadAction<{ message: UserMessage; currentUserId: number }>
    ) {
      const { message, currentUserId } = action.payload;
      console.log(message)
    
      const userId =
        message.MessageSender === message.MessageReceiver
          ? message.MessageSender
          : message.MessageSender === currentUserId
            ? message.MessageReceiver
            : message.MessageSender;
    
      const previousMessages = state.groupedMessages[userId] || [];
      console.log(previousMessages)
      return {
        ...state,
        groupedMessages: {
          ...state.groupedMessages,
          [userId]: [...previousMessages, message],
        },
      };
    },
    deleteConversationGroupedMessage(
      state,
      action: PayloadAction<{ message: any }>
    ) {
      const { userId } = action.payload.message;

      const newGroupedMessages = { ...state.groupedMessages };
      delete newGroupedMessages[userId];
    
      return {
        ...state,
        groupedMessages: newGroupedMessages,
      };
    },
  },
  extraReducers: (builder) => {
    builder.addCase(getAllUserMessages.pending, (state) => {
      state.isLoading = true;
      state.error = null;
    });
    builder.addCase(getAllUserMessages.fulfilled, (state, action) => {
      const messages: UserMessage[] = action.payload.messages;

      state.groupedMessages = messages;
      state.isLoading = false;
    });
    builder.addCase(getAllUserMessages.rejected, (state, action) => {
      state.isLoading = false;
      state.error = action.payload as string;
    });
  }
});

// export const { editMessage, deleteMessage } = userMessagesSlice.actions;
export default userMessagesSlice.reducer;
export const { addGroupedMessage, deleteConversationGroupedMessage } = userMessagesSlice.actions;