<script lang="ts">
    type Message = {
        id: string;
        conversation_id: string;
        sender_id: string;
        content: string;
        media_type: "text" | "image" | "gif" | "video";
        media_url: string;
        created_at: string; // ISO string
        updated_at: string;
        is_deleted: boolean;
        is_unsent: boolean;
    };

  let messages: Message[] = [];
//   const socket = new WebSocket("ws://localhost:50056/ws");
const userId = "ae2cf219-a6e3-46d7-aee0-6f05eb23f602";
const conversationId = "034ff50e-d68a-4ba9-9991-4bd7992ff407";

const socket = new WebSocket(
  `ws://localhost:50056/ws?user_id=${userId}&conversation_id=${conversationId}`
);

  socket.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    console.log(msg);
    messages = [...messages, msg]; // refresh the UI
  };

  socket.onopen = () => {
    console.log("WebSocket connection established");
  };

  socket.onerror = (err) => {
    console.error("WebSocket error", err);
  };
</script>

<div class="home-container">
    <ul>
    {#each messages as msg (msg.id)}
        <li><strong>{msg.sender_id}:</strong> {msg.content}</li>
    {/each}
    </ul>
</div>

<style>
    .home-container {
        display: flex;
        width: 100%;
        height: 100vh;
        background-color: var(--background);
        color: var(--text-primary);
    }
</style>