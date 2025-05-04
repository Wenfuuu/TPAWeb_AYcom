<script lang="ts">
  import { onMount } from "svelte";
  import { type User, type UserThread } from "../model/user";
  import {
    fetchBookmarkWithUser,
    type Thread,
  } from "../model/thread";
  import { navigate } from "svelte-routing";

  import ThreadFeed from "./ThreadFeed.svelte";
  import { debounce } from "../utils/utils";

  export let threads: (Thread & { userData: UserThread })[] = [];
  export let refreshData: (page: number, search: string) => Promise<void>;
  export let user: User | null;

  let page = 1;
  let isLoading = false;
  let hasMore = true;
  let searchQuery = "";

  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;

  const debouncedSearch = debounce(handleSearch, 300);
  async function handleSearch() {
    // search here
    await refreshData(1, searchQuery);
  }

  function observeScroll() {
    if (observer) {
      observer.disconnect();
    }

    if (!contentAreaElement) {
      // Retry after the current event loop tick, allowing bind:this to complete
      requestAnimationFrame(observeScroll);
      return;
    }

    observer = new IntersectionObserver(
      async (entries) => {
        if (entries[0].isIntersecting && !isLoading && hasMore) {
          isLoading = true;

          await new Promise(resolve => setTimeout(resolve, 500));

          let newThreads;
          if(user){
            newThreads = await fetchBookmarkWithUser(page+1, user.id, searchQuery);
          }

          if(!newThreads){
            hasMore = false;
            isLoading = false;
            return;
          }
          threads = [...threads, ...newThreads];
          page += 1;

          isLoading = false;
        }
      },
      {
        root: contentAreaElement,
        rootMargin: "1px",
        threshold: 1,
      }
    );

    if (sentinel) {
      observer.observe(sentinel);
    }
  }

  onMount(async () => {
    observeScroll();
    await refreshData(page, searchQuery);
  });
</script>

{#if user}
    <main class="content-area" bind:this={contentAreaElement}>
        <div class="header">
            <div class="header-top">
                <button type="button" class="back-button" aria-label="Back button" on:click={() => navigate("/home")}>
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="m15 18-6-6 6-6"></path>
                    </svg>
                </button>
                <h1>Bookmarks</h1>
            </div>
            <div class="search-container">
                <div class="search-input-wrapper">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="search-icon">
                        <circle cx="11" cy="11" r="8"></circle>
                        <path d="m21 21-4.3-4.3"></path>
                    </svg>
                    <input
                        type="text"
                        class="search-input"
                        placeholder="Search Bookmarks"
                        bind:value={searchQuery}
                        on:input={debouncedSearch}
                    />
                </div>
            </div>
        </div>

        <ThreadFeed {threads} {user} />

        <div bind:this={sentinel} class="loading-sentinel">
        {#if isLoading}
            <!-- Skeleton Loading UI -->
            <div class="skeleton-loading">
            <!-- Repeat skeleton post structure multiple times for better UX -->
            {#each Array(3) as _, i (i)}
                <div class="skeleton-post">
                <div class="skeleton-avatar pulse"></div>
                <div class="skeleton-content">
                    <div class="skeleton-header">
                    <div class="skeleton-name pulse"></div>
                    <div class="skeleton-username pulse"></div>
                    </div>
                    <div class="skeleton-text pulse"></div>
                    <div class="skeleton-text-short pulse"></div>
                    <!-- Random image placeholder in some posts -->
                    {#if i % 2 === 0}
                    <div class="skeleton-image pulse"></div>
                    {/if}
                    <div class="skeleton-actions">
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    <div class="skeleton-action pulse"></div>
                    </div>
                </div>
                </div>
            {/each}
            </div>
        {:else if !hasMore && threads.length > 0}
            <p>No more thread available!</p>
        {:else if threads.length == 0}
            <p>No bookmark available!</p>
        {/if}
        </div>
    </main>
{/if}

<style>
    .content-area {
        flex: 1;
        border-right: 1px solid var(--border-color);
        width: 50%;
        height: 100vh;
        margin: 0 auto; /* Center if using max-width */
        overflow-y: auto;
    }

    /* Header styles */
    .header {
        position: sticky;
        top: 0;
        background-color: var(--card-background);
        backdrop-filter: blur(12px);
        z-index: 10;
        border-bottom: 1px solid var(--border-color);
        padding: 8px 16px;
    }

    .header-top {
        display: flex;
        align-items: center;
        margin-bottom: 12px;
    }

    .back-button {
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
        margin-right: 20px;
    }

    .back-button:hover {
        background-color: var(--hover-background);
    }

    .header h1 {
        font-size: 20px;
        font-weight: 700;
        color: var(--text-primary);
        margin: 0;
    }

    .search-container {
        margin-bottom: 8px;
    }

    .search-input-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        background-color: var(--hover-background);
        border-radius: 9999px;
        padding: 0 12px;
    }

    .search-icon {
        color: var(--text-tertiary);
        margin-right: 8px;
    }

    .search-input {
        background: transparent;
        border: none;
        color: var(--text-primary);
        font-size: 15px;
        padding: 12px 0;
        width: 100%;
        outline: none;
    }

    .search-input::placeholder {
        color: var(--text-tertiary);
    }

    .loading-sentinel {
        min-height: 100px;
        padding: 20px;
        text-align: center;
        color: var(--text-secondary);
    }

    /* Skeleton loading styles */
    .skeleton-loading {
        width: 90%;
        padding: 0 16px;
    }

    .skeleton-post {
        display: flex;
        padding: 12px 0;
        border-bottom: 1px solid var(--border-color);
        animation: fadeIn 0.5s ease-in-out;
    }

    .skeleton-avatar {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background-color: var(--hover-background);
        margin-right: 12px;
        flex-shrink: 0;
    }

    .skeleton-content {
        flex: 1;
    }

    .skeleton-header {
        display: flex;
        align-items: center;
        margin-bottom: 8px;
        gap: 8px;
    }

    .skeleton-name {
        height: 16px;
        width: 120px;
        background-color: var(--hover-background);
        border-radius: 4px;
    }

    .skeleton-username {
        height: 14px;
        width: 80px;
        background-color: var(--hover-background);
        border-radius: 4px;
    }

    .skeleton-text {
        height: 14px;
        width: 100%;
        background-color: var(--hover-background);
        border-radius: 4px;
        margin-bottom: 8px;
    }

    .skeleton-text-short {
        height: 14px;
        width: 70%;
        background-color: var(--hover-background);
        border-radius: 4px;
        margin-bottom: 16px;
    }

    .skeleton-image {
        height: 180px;
        width: 100%;
        background-color: var(--hover-background);
        border-radius: 16px;
        margin-bottom: 16px;
    }

    .skeleton-actions {
        display: flex;
        justify-content: space-between;
        max-width: 425px;
    }

    .skeleton-action {
        height: 18px;
        width: 36px;
        background-color: var(--hover-background);
        border-radius: 4px;
    }

    .pulse {
        animation: pulse 1.5s infinite;
    }

    @keyframes pulse {
        0% {
            opacity: 0.6;
        }
        50% {
            opacity: 1;
        }
        100% {
            opacity: 0.6;
        }
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }

    @media (max-width: 768px) {
        .content-area {
            border-right: none; /* Remove border on mobile */
            width: 100%;
        }
    }
</style>