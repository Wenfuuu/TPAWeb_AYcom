<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import MicrophoneInput from "./MicrophoneInput.svelte";

  export let placeholder = "Search...";
  export let value = "";

  const dispatch = createEventDispatcher();

  function handleSubmit() {
    dispatch("search", { query: value });
  }

  function handleInput(event: Event) {
    const inputEvent = event as InputEvent;
    const target = inputEvent.target as HTMLInputElement;
    value = target.value;
    dispatch("input", { value });
  }

  function handleVoiceInput(event: CustomEvent<{ text: string }>) {
    value = event.detail.text;
    dispatch("input", { value });

    if (value.trim()) {
      handleSubmit();
    }
  }
</script>

<div class="search-container">
  <form on:submit|preventDefault={handleSubmit} class="search-form">
    <input
      type="text"
      class="search-input"
      {placeholder}
      bind:value
      on:input={handleInput}
    />

    <div class="search-actions">
      <MicrophoneInput
        on:transcribed={handleVoiceInput}
        on:error={(e) => dispatch("error", e.detail)}
        buttonClass="search-mic-button"
      />

      <button type="submit" class="search-button" aria-label="search button">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="feather feather-search"
        >
          <circle cx="11" cy="11" r="8"></circle>
          <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
        </svg>
      </button>
    </div>
  </form>
</div>

<style>
  .search-container {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
  }

  .search-form {
    display: flex;
    position: relative;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    border-radius: 24px;
    overflow: hidden;
  }

  .search-input {
    flex: 1;
    padding: 12px 16px;
    border: 1px solid #ddd;
    border-radius: 24px;
    font-size: 16px;
    outline: none;
    transition: border-color 0.2s;
  }

  .search-input:focus {
    border-color: #2196f3;
  }

  .search-actions {
    position: absolute;
    right: 0;
    top: 0;
    height: 100%;
    display: flex;
    align-items: center;
  }

  .search-button {
    background: none;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 16px;
    height: 100%;
    color: #555;
  }

  .search-button:hover {
    color: #2196f3;
  }

  .feather-search {
    width: 20px;
    height: 20px;
  }
</style>