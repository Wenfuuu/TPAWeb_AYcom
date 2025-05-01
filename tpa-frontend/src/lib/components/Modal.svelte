<!-- src/components/Modal.svelte -->
<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from "svelte";
  import { fade, fly } from "svelte/transition";

  export let show: boolean = false; // Controls visibility

  const dispatch = createEventDispatcher();

  function close() {
    dispatch("close");
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === "Escape") {
      close();
    }
  }

  function handleOuterClick(event: MouseEvent) {
    // Close if clicking directly on the overlay, not its children
    if (event.target === event.currentTarget) {
        close();
    }
  }

  onMount(() => {
    window.addEventListener("keydown", handleKeydown);
  });

  onDestroy(() => {
    window.removeEventListener("keydown", handleKeydown);
  });
</script>

{#if show}
  <!-- svelte-ignore a11y-click-events-have-key-events a11y-no-static-element-interactions -->
  <div
    class="modal-overlay"
    on:click={handleOuterClick}
    transition:fade={{ duration: 150 }}
  >
    <div class="modal-content" role="dialog" aria-modal="true" transition:fly={{ y: 50, duration: 200 }}>
      <button type="button" class="modal-close-button" on:click={close} aria-label="Close modal">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
      <!-- Content goes here -->
      <slot />
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.6); /* Darker overlay */
    display: flex;
    justify-content: center;
    align-items: flex-start; /* Align to top to allow scrolling */
    z-index: 1000;
    padding-top: 5vh; /* Add some top padding */
    overflow-y: auto; /* Allow overlay to scroll if content is too tall */
  }

  .modal-content {
    background-color: var(--card-background);
    padding: 0; /* Remove padding, ComposePost has its own */
    border-radius: 16px;
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
    position: relative;
    width: 90%;
    max-width: 600px; /* Max width for the compose modal */
    z-index: 1001;
    /* Remove fixed height, let content dictate */
    /* max-height: 90vh; */ /* Set max height */
    /* overflow-y: auto; */ /* Allow content itself to scroll if needed */
    border: 1px solid var(--border-color);
  }

  .modal-close-button {
    position: absolute;
    top: 10px;
    left: 10px; /* Position top-left */
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: var(--text-tertiary);
    padding: 8px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1002;
  }
  .modal-close-button:hover {
      background-color: var(--hover-background);
      color: var(--text-primary);
  }

   @media (max-width: 768px) {
     .modal-content {
         width: 95%;
         max-width: none;
         margin-top: 2vh; /* Smaller top margin */
     }
     .modal-overlay {
         padding-top: 2vh;
     }
   }

</style>