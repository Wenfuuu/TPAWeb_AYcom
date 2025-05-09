<script lang="ts">
  import { onMount } from "svelte";
  import {
  createNewConversation,
    fetchMessageWithUser,
    fetchUserConversation,
    updateConversationParticipants,
    type Conversation,
    type Message,
  } from "../lib/model/chat";
  import { checkAuth, getAuthUser } from "../lib/utils/auth";
  import { navigate } from "svelte-routing";
  import { toastStore } from "../lib/stores/toast";
  import LeftSidebar from "../lib/components/LeftSidebar.svelte";
  import { fetchUser, fetchUserEmailById, type User } from "../lib/model/user";
  import { uploadImage, uploadVideo } from "../lib/utils/supabaseStorage";

  let isDarkMode = true;
  let user: User;
  let allUsers: User[] = [];

  let conversations: Conversation[] = [];
  let currentConversation: Conversation | null = null;
  let messages: Message[] = [];
  let newMessage = "";
  let searchQuery = "";

  let uploadedImage = "";
  let imageFile: File | null = null;
  let supabaseURL: string = "";

  let videoFile: File | null = null;
  let videoPreviewUrl = "";
  let videoSupabaseURL = "";
  let hasVideo = false;

  let socket: WebSocket;
  let showParticipantsDropdown = false;
  let selectedParticipants: string[] = [];

  let showUserDropdown = false;
  let selectedUsers: string[] = [];

  function setupWebSocket(userId: string, conversationId: string) {
    if (socket) {
      socket.close();
    }

    socket = new WebSocket(
      `ws://localhost:50056/ws?user_id=${userId}&conversation_id=${conversationId}`
    );

    socket.onmessage = (event) => {
        const wsMessage = JSON.parse(event.data);
        // console.log("received wsmssage: " + wsMessage);
        switch (wsMessage.type) {
            case "chat": {
                console.log("received a message");
                // Handle chat message
                const chatMessage = wsMessage.payload;
                messages = [...messages, chatMessage];
                // Your existing chat message handling code
                break;
            }

            case "notification": {
                console.log("received notification");
                // Handle notification
                const notification = wsMessage.payload;
                toastStore.info(notification.content);
                break;
            }

            default:
                console.log("Unknown message type:", wsMessage.type);
        }
    };

    socket.onopen = () => {
      console.log("WebSocket connection established");
    };

    socket.onerror = (err) => {
      if (err instanceof Error) {
        toastStore.error(err.message);
      } else {
        toastStore.error("An unexpected error occurred!");
      }
    };
  }

  async function createConversation() {
    console.log("creating new conversation");
    console.log(selectedUsers);

    await createNewConversation(selectedUsers, user.id);
    showUserDropdown = false;
    toastStore.success("Success creating a new conversation!");
    await refreshData();
  }

  async function updateParticipants() {
    if(!currentConversation) return;
    console.log(selectedParticipants);
    await updateConversationParticipants(currentConversation.id, selectedParticipants);
    showParticipantsDropdown = false;
    toastStore.success("Success update participants!");
    await refreshData();
  }

  async function selectConversation(conversation: Conversation) {
    currentConversation = conversation;
    console.log(currentConversation.participants);

    if (currentConversation && currentConversation.participants) {
      selectedParticipants = currentConversation.participants.map(p => p.user_id);
    } else {
      selectedParticipants = [];
    }

    await refreshData();
  }

  async function uploadImagesToSupabase(): Promise<boolean> {
    if (!imageFile) return false;

    const supaURL = await uploadImage(imageFile);
    if (supaURL) supabaseURL = supaURL;

    return true;
  }

  function handleMediaUpload(event: Event): void {
    const target = event.target as HTMLInputElement;
    const files = target.files;

    if (!files || files.length === 0) return;

    for (let i = 0; i < files.length; i++) {
      const file = files[i];

      // Check if the file is a video
      if (file.type.startsWith("video/")) {
        // Only allow one video
        if (hasVideo) {
          alert(
            "Only one video allowed. Please remove the existing video first."
          );
          continue;
        }

        videoFile = file;
        videoPreviewUrl = URL.createObjectURL(file);
        hasVideo = true;
      } else if (file.type.startsWith("image/")) {
        // Handle image files
        imageFile = file;
        const reader = new FileReader();

        reader.onload = (e: ProgressEvent<FileReader>) => {
          const result = e.target?.result;
          if (typeof result === "string") {
            uploadedImage = result;
          }
        };

        reader.readAsDataURL(file);
      }
    }
  }

  async function refreshData() {
    conversations = await fetchUserConversation(user.id);

    if (!currentConversation) return;

    // establish websocket
    setupWebSocket(user.id, currentConversation.id);

    messages = await fetchMessageWithUser(currentConversation.id);
  }

  function toggleTheme() {
    isDarkMode = !isDarkMode;
    const newTheme = isDarkMode ? "dark" : "light";
    localStorage.setItem("theme", newTheme);
    document.documentElement.setAttribute("data-theme", newTheme);
  }

  async function sendMessage() {
    if (
      !newMessage.trim() ||
      !socket ||
      socket.readyState !== WebSocket.OPEN ||
      !currentConversation
    )
      return;

    let message;
    if (imageFile) {
      await uploadImagesToSupabase();
      message = {
        conversation_id: currentConversation.id,
        sender_id: user.id,
        content: newMessage,
        media_type: "image",
        media_url: supabaseURL,
      };
    } else if (videoFile) {
      const videoUrl = await uploadVideo(videoFile);
      if (videoUrl) {
        videoSupabaseURL = videoUrl;
      }
      message = {
        conversation_id: currentConversation.id,
        sender_id: user.id,
        content: newMessage,
        media_type: "video",
        media_url: videoSupabaseURL,
      };
    } else {
      message = {
        conversation_id: currentConversation.id,
        sender_id: user.id,
        content: newMessage,
        media_type: "text",
        media_url: "",
      };
    }

    const socketMessage = {
        type: "chat",
        payload: message
    };

    socket.send(JSON.stringify(socketMessage));

    // test notification
    const otherParticipants = currentConversation.participants?.filter(p => p.user_id !== user.id) || [];
    // For each participant, send a notification
    for (const participant of otherParticipants) {
        const email = await fetchUserEmailById(participant.user_id, user.id);
        console.log(email);

        const notif = {
            recipient_id: participant.user_id,
            sender_id: user.id,
            content: `New message from ${user.name} in ${currentConversation?.name}`,
            category: "message",
            email: email
        };

        const notifMessage = {
            type: "notification",
            payload: notif
        };

        socket.send(JSON.stringify(notifMessage));
    }

    newMessage = "";
    uploadedImage = "";
    imageFile = null;
    supabaseURL = "";
    videoFile = null;
    videoPreviewUrl = "";
    videoSupabaseURL = "";
    hasVideo = false;
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === "Enter" && !event.shiftKey) {
      event.preventDefault();
      sendMessage();
    }
  }

  function formatTime(timestamp: string): string {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
  }

  function formatDate(timestamp: string): string {
    const date = new Date(timestamp);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === today.toDateString()) {
      return "Today";
    } else if (date.toDateString() === yesterday.toDateString()) {
      return "Yesterday";
    } else {
      return date.toLocaleDateString([], {
        month: "short",
        day: "numeric",
        year: "numeric",
      });
    }
  }

  function filterConversations() {
    if (!searchQuery) return conversations;
    return conversations.filter((conv) =>
      conv.name.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }

  function getLastMessage(conversationId: string) {
    // This would typically come from the API, but for now we'll use the messages array
    const conversationMessages = messages.filter(
      (m) => m.conversation_id === conversationId
    );
    if (conversationMessages.length > 0) {
      const lastMsg = conversationMessages[conversationMessages.length - 1];
      return lastMsg.content;
    }
    return "";
  }

  function getLastMessageTime(conversationId: string) {
    const conversationMessages = messages.filter(
      (m) => m.conversation_id === conversationId
    );
    if (conversationMessages.length > 0) {
      const lastMsg = conversationMessages[conversationMessages.length - 1];
      const date = new Date(lastMsg.created_at);
      const now = new Date();

      // If today, return time
      if (date.toDateString() === now.toDateString()) {
        return formatTime(lastMsg.created_at);
      }

      // If this year, return month and day
      if (date.getFullYear() === now.getFullYear()) {
        return date.toLocaleDateString([], { month: "short", day: "numeric" });
      }

      // Otherwise return full date
      return date.toLocaleDateString([], {
        month: "short",
        day: "numeric",
        year: "numeric",
      });
    }
    return "";
  }

  onMount(async () => {
    await checkAuth(navigate);

    const usr = await getAuthUser();
    if (usr) user = usr;

    allUsers = await fetchUser();
    // console.log(allUsers);
    // conversations = await fetchUserConversation(user.id);
    // console.log(conversations);

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme) {
      isDarkMode = savedTheme === "dark";
      document.documentElement.setAttribute("data-theme", savedTheme);
    } else {
      const prefersDark = window.matchMedia(
        "(prefers-color-scheme: dark)"
      ).matches;
      isDarkMode = prefersDark;
      document.documentElement.setAttribute(
        "data-theme",
        prefersDark ? "dark" : "light"
      );
    }

    await refreshData();
  });
</script>

{#if conversations || user}
  <div class="home-container">
    <LeftSidebar {toggleTheme} {isDarkMode} {user} {refreshData} />

    <div class="messages-container">
      <div class="conversations-list">
        <div class="conversations-header">
          <h2>Messages</h2>
          <div class="header-actions">
            <button type="button" aria-label="icon button" class="icon-button"
            on:click={() => showUserDropdown = !showUserDropdown}>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path
                  d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                ></path>
              </svg>
            </button>

            {#if showUserDropdown}
            <div class="create-dropdown-menu">
                <div class="create-dropdown-header">Start New Conversation</div>
                {#each allUsers as user (user.id)}
                <label class="create-dropdown-item">
                    <input type="checkbox" bind:group={selectedUsers} value={user.id} />
                    <div class="participant-avatar">
                    <img src={user.profile_picture || "/placeholder.svg"} alt={user.name} />
                    </div>
                    <span>{user.name}</span>
                </label>
                {/each}
                {#if selectedUsers.length > 0}
                <button type="button" class="create-convo-btn" on:click={async() => {
                    await createConversation();
                }}>
                    Create Conversation
                </button>
                {/if}
            </div>
            {/if}
          </div>
        </div>

        <div class="search-container">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
          </svg>
          <input
            type="text"
            placeholder="Search Direct Messages"
            bind:value={searchQuery}
          />
        </div>

        <div class="conversations">
          {#if filterConversations().length === 0}
            <div class="empty-conversations">
              <p>No conversations found</p>
            </div>
          {:else}
            {#each filterConversations() as conversation (conversation.id)}
              <div
                role="button"
                tabindex="0"
                class="conversation-item"
                class:active={currentConversation &&
                  currentConversation.id === conversation.id}
                on:click={() => selectConversation(conversation)}
                on:keydown={(event) => {
                  if (event.key === "Enter" || event.key === " ") {
                    selectConversation(conversation);
                  }
                }}
              >
                <div class="conversation-details">
                  <div class="conversation-header">
                    <span class="conversation-name">{conversation.name}</span>
                    <span class="conversation-time"
                      >· {getLastMessageTime(conversation.id) || "10h"}</span
                    >
                  </div>
                  <div class="conversation-preview">
                    {getLastMessage(conversation.id) || "No messages yet"}
                  </div>
                </div>
              </div>
            {/each}
          {/if}
        </div>
      </div>

      <div class="chat-detail">
        {#if currentConversation}
          <div class="chat-header">
            <div class="chat-header-info">
              <span class="chat-name">{currentConversation.name}</span>
              <button
              aria-label="icon button"
              type="button"
              class="icon-button"
              on:click={() =>
                (showParticipantsDropdown = !showParticipantsDropdown)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="10"></circle>
                <path d="M12 16v-4"></path>
                <path d="M12 8h.01"></path>
              </svg>
            </button>
            </div>
          </div>

          <div class="dropdown-container">
            {#if showParticipantsDropdown}
              <div class="dropdown-menu">
                <div class="dropdown-header">Participants</div>
                {#each allUsers as user (user.id)}
                <label class="dropdown-item">
                    <input
                    type="checkbox"
                    bind:group={selectedParticipants}
                    value={user.id}
                    checked={currentConversation?.participants?.some(p => p.user_id === user.id)}
                    />
                    <div class="participant-avatar">
                    <img
                        src={user.profile_picture || "/placeholder.svg?height=24&width=24"}
                        alt={user.name}
                    />
                    </div>
                    <span>{user.username}</span>
                </label>
                {/each}

                {#if selectedParticipants.length > 0}
                <button
                    type="button"
                    class="add-participants-btn"
                    on:click={async() => {
                        await updateParticipants();
                    }}
                >
                    Update Participants
                </button>
                {/if}
              </div>
            {/if}
          </div>

          <div class="chat-profile">
            <div class="profile-name">{currentConversation.name}</div>
            <div class="profile-meta">
              Conversation started {new Date(
                currentConversation.created_at
              ).toLocaleDateString([], {
                month: "long",
                day: "numeric",
                year: "numeric",
              })}
            </div>
          </div>

          <div class="messages-list">
            {#if messages.length === 0}
              <div class="empty-messages">
                <p>No messages yet</p>
                <span>Start a conversation with {currentConversation.name}</span
                >
              </div>
            {:else}
              {#each messages as message, i (message.id)}
                {#if i === 0 || formatDate(messages[i - 1].created_at) !== formatDate(message.created_at)}
                  <div class="date-divider">
                    <span>{formatDate(message.created_at)}</span>
                  </div>
                {/if}
                <div
                  class="message-item"
                  class:own-message={message.sender_id === user.id}
                >
                  {#if message.sender_id !== user.id}
                    <div class="message-avatar">
                      <img
                        src={message.profile_picture ||
                          "/placeholder.svg?height=32&width=32"}
                        alt="User"
                      />
                    </div>
                  {/if}
                  <div class="message-content">
                    {#if message.media_type === "text"}
                      <div class="message-bubble">
                        {message.content}
                      </div>
                    {:else if message.media_type === "image"}
                      <div class="message-media">
                        <img src={message.media_url} alt="" />
                        {#if message.content}
                          <div class="message-bubble">
                            {message.content}
                          </div>
                        {/if}
                      </div>
                    {:else if message.media_type === "gif"}
                      <div class="message-media">
                        <img
                          src={message.media_url || "/placeholder.svg"}
                          alt="GIF"
                        />
                        {#if message.content}
                          <div class="message-bubble">
                            {message.content}
                          </div>
                        {/if}
                      </div>
                    {:else if message.media_type === "video"}
                      <div class="message-media">
                        <video src={message.media_url} controls>
                          <track
                            kind="captions"
                            src="path_to_captions.vtt"
                            srclang="en"
                            label="English"
                          />
                        </video>
                        {#if message.content}
                          <div class="message-bubble">
                            {message.content}
                          </div>
                        {/if}
                      </div>
                    {/if}
                    <div class="message-time">
                      {formatTime(message.created_at)}
                      {#if message.is_unsent}
                        <span class="unsent-indicator">Unsent</span>
                      {/if}
                    </div>
                  </div>
                </div>
              {/each}
            {/if}
          </div>

          <div class="message-input-container">
            <!-- Media Previews -->
            {#if uploadedImage || videoPreviewUrl}
              <div class="media-preview-container">
                {#if uploadedImage}
                  <div class="media-preview">
                    <img src={uploadedImage} alt="Upload preview" />
                    <button
                      type="button"
                      aria-label="remove button"
                      class="remove-media-button"
                      on:click={() => {
                        uploadedImage = "";
                        imageFile = null;
                      }}
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      >
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                      </svg>
                    </button>
                  </div>
                {/if}

                {#if videoPreviewUrl}
                  <div class="media-preview">
                    <video src={videoPreviewUrl} controls>
                      <track
                        kind="captions"
                        src=""
                        srclang="en"
                        label="English"
                      />
                    </video>
                    <button
                      type="button"
                      aria-label="remove button"
                      class="remove-media-button"
                      on:click={() => {
                        videoPreviewUrl = "";
                        videoFile = null;
                        hasVideo = false;
                      }}
                    >
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      >
                        <line x1="18" y1="6" x2="6" y2="18"></line>
                        <line x1="6" y1="6" x2="18" y2="18"></line>
                      </svg>
                    </button>
                  </div>
                {/if}
              </div>
            {/if}

            <div class="message-input-wrapper">
              <div class="message-input-actions">
                <label class="action-button" for="media-upload">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <rect width="18" height="18" x="3" y="3" rx="2" ry="2"
                    ></rect>
                    <circle cx="9" cy="9" r="2"></circle>
                    <path d="m21 15-3.086-3.086a2 2 0 0 0-2.828 0L6 21"></path>
                  </svg>
                  <input
                    type="file"
                    id="media-upload"
                    accept="image/*,video/*"
                    on:change={handleMediaUpload}
                    class="hidden-input"
                  />
                </label>
              </div>
              <textarea
                placeholder="Start a new message"
                bind:value={newMessage}
                on:keydown={handleKeyDown}
              ></textarea>
              <button
                type="button"
                aria-label="send button"
                class="send-button"
                class:active={newMessage.trim().length > 0}
                on:click={sendMessage}
                disabled={!newMessage.trim()}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="m22 2-7 20-4-9-9-4Z"></path>
                  <path d="M22 2 11 13"></path>
                </svg>
              </button>
            </div>
          </div>
        {:else}
          <div class="no-conversation-selected">
            <div class="no-conversation-content">
              <h2>Select a conversation</h2>
              <p>Choose from your existing conversations or start a new one</p>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .home-container {
    display: flex;
    width: 100%;
    height: 100vh;
    background-color: var(--background);
    color: var(--text-primary);
  }

  .messages-container {
    display: flex;
    flex: 1;
    height: 100%;
  }

  /* Conversations List Styles */
  .conversations-list {
    width: 380px;
    border-right: 1px solid var(--border-color);
    display: flex;
    flex-direction: column;
    height: 100%;
  }

  .conversations-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .conversations-header h2 {
    font-size: 20px;
    font-weight: 700;
    margin: 0;
  }

  .header-actions {
    display: flex;
    gap: 8px;
    position: relative;
  }

  .create-dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: var(--background, #fff);
  border: 1px solid var(--border-color, #ddd);
  border-radius: 8px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  padding: 12px;
  min-width: 240px;
  z-index: 100;
  margin-top: 10px;
}

.create-dropdown-header {
  font-weight: bold;
  font-size: 14px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color, #eee);
  margin-bottom: 8px;
}

.create-dropdown-item {
  display: flex;
  align-items: center;
  padding: 6px 0;
  cursor: pointer;
  transition: background-color 0.2s;
  font-size: 14px;
}

.create-dropdown-item:hover {
  background-color: var(--hover-background, #f6f6f6);
  border-radius: 4px;
}

.create-dropdown-item input[type="checkbox"] {
  margin-right: 8px;
  transform: scale(1.1);
}

  .icon-button {
    background: none;
    border: none;
    color: var(--text-primary);
    cursor: pointer;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .icon-button:hover {
    background-color: var(--hover-background);
  }

  .search-container {
    display: flex;
    align-items: center;
    background-color: var(--hover-background);
    border-radius: 9999px;
    margin: 12px 16px;
    padding: 0 12px;
  }

  .search-container svg {
    color: var(--text-tertiary);
    margin-right: 8px;
  }

  .search-container input {
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 15px;
    padding: 10px 0;
    width: 100%;
    outline: none;
  }

  .search-container input::placeholder {
    color: var(--text-tertiary);
  }

  .conversations {
    flex: 1;
    overflow-y: auto;
  }

  .conversation-item {
    display: flex;
    padding: 12px 16px;
    cursor: pointer;
    transition: background-color 0.2s;
    border-bottom: 1px solid var(--border-color);
  }

  .conversation-item:hover {
    background-color: var(--hover-background);
  }

  .conversation-item.active {
    background-color: var(--hover-background);
  }

  .conversation-details {
    flex: 1;
    min-width: 0;
  }

  .conversation-header {
    display: flex;
    align-items: center;
    margin-bottom: 4px;
  }

  .conversation-name {
    font-weight: 700;
    color: var(--text-primary);
    margin-right: 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .conversation-time {
    color: var(--text-tertiary);
    font-size: 14px;
    white-space: nowrap;
  }

  .conversation-preview {
    color: var(--text-tertiary);
    font-size: 14px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .empty-conversations {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: var(--text-tertiary);
  }

  /* Chat Detail Styles */
  .chat-detail {
    flex: 1;
    display: flex;
    flex-direction: column;
    height: 100%;
    position: relative;
  }

  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .chat-header-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
  }

  .chat-name {
    font-size: 18px;
    font-weight: 700;
    margin-right: 4px;
  }

  .chat-profile {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 24px 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .profile-name {
    font-size: 18px;
    font-weight: 700;
    margin-bottom: 12px;
  }

  .profile-meta {
    text-align: center;
    color: var(--text-tertiary);
    font-size: 14px;
  }

  .messages-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    display: flex;
    flex-direction: column;
  }

  .date-divider {
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 16px 0;
    position: relative;
  }

  .date-divider span {
    background-color: var(--background);
    padding: 0 12px;
    color: var(--text-tertiary);
    font-size: 14px;
    position: relative;
    z-index: 1;
  }

  .date-divider::before {
    content: "";
    position: absolute;
    left: 0;
    right: 0;
    height: 1px;
    background-color: var(--border-color);
  }

  .message-item {
    display: flex;
    margin-bottom: 16px;
    max-width: 80%;
  }

  .own-message {
    align-self: flex-end;
    flex-direction: row-reverse;
  }

  .message-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    overflow: hidden;
    margin-right: 8px;
  }

  .own-message .message-avatar {
    margin-right: 0;
    margin-left: 8px;
  }

  .message-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .message-content {
    display: flex;
    flex-direction: column;
  }

  .own-message .message-content {
    align-items: flex-end;
  }

  .message-bubble {
    background-color: var(--hover-background);
    border-radius: 18px;
    padding: 12px 16px;
    color: var(--text-primary);
    max-width: 100%;
    word-wrap: break-word;
  }

  .own-message .message-bubble {
    background-color: var(--button-primary);
    color: white;
  }

  .message-media {
    max-width: 300px;
    border-radius: 16px;
    overflow: hidden;
    margin-bottom: 4px;
  }

  .message-media img,
  .message-media video {
    width: 100%;
    max-height: 300px;
    object-fit: contain;
    background-color: var(--hover-background);
    border-radius: 16px;
  }

  .message-time {
    font-size: 12px;
    color: var(--text-tertiary);
    margin-top: 4px;
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .unsent-indicator {
    font-style: italic;
    color: var(--text-tertiary);
  }

  .empty-messages {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 200px;
    color: var(--text-tertiary);
    text-align: center;
  }

  .empty-messages p {
    font-weight: 700;
    margin-bottom: 8px;
  }

  .media-preview-container {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 8px;
    padding: 0 16px;
  }

  .media-preview {
    position: relative;
    width: 120px;
    height: 120px;
    border-radius: 12px;
    overflow: hidden;
    background-color: var(--hover-background);
  }

  .media-preview img,
  .media-preview video {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .remove-media-button {
    position: absolute;
    top: 4px;
    right: 4px;
    background-color: rgba(0, 0, 0, 0.6);
    color: white;
    border: none;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0.8;
    transition: opacity 0.2s;
  }

  .remove-media-button:hover {
    opacity: 1;
  }

  .message-input-container {
    padding: 16px;
    border-top: 1px solid var(--border-color);
  }

  .message-input-wrapper {
    display: flex;
    align-items: center;
    background-color: var(--hover-background);
    border-radius: 24px;
    padding: 8px 12px;
  }

  .message-input-actions {
    display: flex;
    gap: 8px;
  }

  .action-button {
    background: none;
    border: none;
    color: var(--button-primary);
    cursor: pointer;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s;
  }

  .action-button:hover {
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
  }

  textarea {
    flex: 1;
    background: none;
    border: none;
    color: var(--text-primary);
    font-size: 15px;
    padding: 8px;
    resize: none;
    height: 24px;
    max-height: 120px;
    outline: none;
    font-family: inherit;
  }

  .send-button {
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: not-allowed;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0.5;
  }

  .send-button.active {
    color: var(--button-primary);
    cursor: pointer;
    opacity: 1;
  }

  .no-conversation-selected {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
  }

  .no-conversation-content {
    text-align: center;
    max-width: 300px;
  }

  .no-conversation-content h2 {
    margin-bottom: 8px;
  }

  .no-conversation-content p {
    color: var(--text-tertiary);
  }

  .hidden-input {
    display: none;
  }

  .dropdown-container {
  position: relative;
  display: inline-block;
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: var(--background, #fff);
  border: 1px solid var(--border-color, #ddd);
  border-radius: 12px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.1);
  padding: 8px 0;
  min-width: 220px;
  z-index: 100;
  margin-top: 8px;
  overflow: hidden;
}

.dropdown-header {
  font-weight: 600;
  font-size: 13px;
  color: #555;
  padding: 8px 16px;
  border-bottom: 1px solid var(--border-color, #eee);
  background-color: var(--header-background, #f9f9f9);
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.dropdown-item:hover {
  background-color: var(--hover-background, #f0f0f0);
}

.participant-avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #ddd;
}

.participant-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.dropdown-item input[type="checkbox"] {
  accent-color: var(--accent-color, #007bff);
  margin: 0;
}

.add-participants-btn,
.create-convo-btn {
  width: calc(100% - 32px);
  margin: 12px 16px 8px;
  padding: 8px;
  border: none;
  border-radius: 8px;
  background-color: var(--accent-color, #007bff);
  color: #fff;
  font-weight: 500;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.add-participants-btn:hover {
  background-color: var(--accent-hover, #005fcc);
}

  @media (max-width: 768px) {
    .conversations-list {
      width: 100%;
      display: none;
    }

    .chat-detail {
      width: 100%;
    }
  }
</style>
