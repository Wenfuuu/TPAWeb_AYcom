import { toastStore } from "../stores/toast";

export function setupWS(socket: WebSocket, userId: string) {
  if (socket) {
    socket.close();
  }

  socket = new WebSocket(`ws://localhost:50056/ws?user_id=${userId}`);

  socket.onmessage = (event) => {
    const wsMessage = JSON.parse(event.data);
    // console.log(wsMessage.type);
    switch (wsMessage.type) {
      case "notification": {
        const notification = wsMessage.payload;
        toastStore.info(notification.content);
        break;
      }

      default:
      // console.log("Unknown message type:", wsMessage.type);
    }
  };

  socket.onopen = () => {
    // console.log("WebSocket connection established");
  };

  socket.onerror = (err) => {
    if (err instanceof Error) {
      toastStore.error(err.message);
    } else {
      toastStore.error("An unexpected error occurred!");
    }
  };

  // socket.onclose = () => {
  //   console.log("WebSocket connection closed. Reconnecting...");
  //   setTimeout(() => setupWebSocket(userId), 5000); // Reconnect after 5 seconds
  // };

  return socket;
}

export function sendMessageNotif(socket: WebSocket,
    recipient_id: string, sender_id: string,
    content: string, email: string) {

    const notif = {
        recipient_id,
        sender_id,
        content,
        category: "message",
        email,
    };

    const notifMessage = {
        type: "notification",
        payload: notif,
    };

    socket.send(JSON.stringify(notifMessage));
}

export function sendThreadNotif(socket: WebSocket,
  recipient_id: string, sender_id: string, thread_id: string,
  content: string, email: string) {

  const notif = {
      recipient_id,
      sender_id,
      thread_id,
      content,
      category: "message",
      email,
  };

  const notifMessage = {
      type: "notification",
      payload: notif,
  };

  socket.send(JSON.stringify(notifMessage));
}