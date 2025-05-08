<script lang="ts">
  import { onDestroy } from "svelte";
  import { createEventDispatcher } from "svelte";

  // Props
  export let listening = false;
  export let buttonClass = "";
  export let iconClass = "";

  // Event dispatcher
  const dispatch = createEventDispatcher();

  // Variables
  let mediaRecorder: MediaRecorder;
  let audioChunks: Blob[] = [];
  let stream: MediaStream;

  // Start recording
  async function startListening() {
    try {
      if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        alert("Your browser does not support audio recording");
        return;
      }

      stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      listening = true;

      mediaRecorder = new MediaRecorder(stream);
      audioChunks = [];

      mediaRecorder.addEventListener("dataavailable", (event) => {
        audioChunks.push(event.data);
      });

      mediaRecorder.addEventListener("stop", () => {
        const audioBlob = new Blob(audioChunks, { type: "audio/wav" });
        sendAudioToServer(audioBlob);
      });

      mediaRecorder.start();
      dispatch("start");
    } catch (err) {
      console.error("Error accessing microphone:", err);
    }
  }

  // Stop recording
  function stopListening() {
    if (mediaRecorder && mediaRecorder.state !== "inactive") {
      mediaRecorder.stop();
      listening = false;

      // Stop all tracks
      if (stream) {
        stream.getTracks().forEach((track) => track.stop());
      }

      dispatch("stop");
    }
  }

  // Toggle recording
  function toggleListening() {
    if (listening) {
      stopListening();
    } else {
      startListening();
    }
  }

  // Send audio to backend
  async function sendAudioToServer(audioBlob: Blob): Promise<void> {
    try {
      // Create form data
      const formData = new FormData();
      formData.append("file", audioBlob);

      // Send to AI service
      const response = await fetch("http://localhost:5000/transcribe", {
        method: "POST",
        body: formData,
      });

      const result: { success: boolean; text?: string; error?: string } =
        await response.json();

      if (result.success && result.text) {
        // Dispatch the transcribed text
        dispatch("transcribed", { text: result.text });
      } else {
        console.error("Transcription failed:", result.error);
        dispatch("error", { message: result.error || "Transcription failed" });
      }
    } catch (err) {
      console.error("Error sending audio to server:", err);
      dispatch("error", {
        message: "Error sending audio: " + (err as Error).message,
      });
    }
  }

  // Clean up on component destroy
  onDestroy(() => {
    if (listening) {
      stopListening();
    }
  });
</script>

<button
  type="button"
  on:click={toggleListening}
  class={`microphone-button ${buttonClass} ${listening ? "listening" : ""}`}
  aria-label={listening ? "Stop listening" : "Start listening"}
  title={listening ? "Stop listening" : "Start listening"}
>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
    class={`feather feather-mic ${iconClass}`}
  >
    {#if listening}
      <!-- Mic with waves (recording) -->
      <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"></path>
      <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
      <line x1="12" y1="19" x2="12" y2="23"></line>
      <line x1="8" y1="23" x2="16" y2="23"></line>
      <path d="M4 10a8 8 0 0 0 16 0" stroke-dasharray="1, 3"></path>
    {:else}
      <!-- Normal mic -->
      <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"></path>
      <path d="M19 10v2a7 7 0 0 1-14 0v-2"></path>
      <line x1="12" y1="19" x2="12" y2="23"></line>
      <line x1="8" y1="23" x2="16" y2="23"></line>
    {/if}
  </svg>
</button>

<style>
  .microphone-button {
    background: none;
    border: none;
    cursor: pointer;
    padding: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all 0.2s ease;
  }

  .microphone-button:hover {
    background-color: rgba(0, 0, 0, 0.05);
  }

  .microphone-button.listening {
    background-color: rgba(255, 0, 0, 0.1);
  }

  .feather-mic {
    width: 24px;
    height: 24px;
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
    }
    50% {
      transform: scale(1.1);
    }
    100% {
      transform: scale(1);
    }
  }

  .listening .feather-mic {
    color: red;
    animation: pulse 1.5s infinite;
  }
</style>
