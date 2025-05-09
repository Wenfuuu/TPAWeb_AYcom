<script lang="ts">
  import { onMount } from "svelte";
  import { fetchLatestThread, type Thread } from "../model/thread";
  import type { User, UserThread } from "../model/user";
  import SearchBar from "./SearchBar.svelte";
  import ThreadFeed from "./ThreadFeed.svelte";
  import type { Community } from "../model/community";

    export let user: User | null;
    export let refreshData: (page: number,
    currTab: "top" | "latest" | "people" | "media" | "communities",
    search: string,) => Promise<void>;
    export let latestThreads: (Thread & { userData: UserThread })[] = [];
    export let communities: Community[] = [];


  let activeTab: "top" | "latest" | "people" | "media" | "communities";
  let page = 1;

  let hasMore = true;
  let searchQuery: string = "";
  let isLoading: boolean = false;
  let errorMessage: string = "";

  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;

  async function handleSearch(
    event: CustomEvent<{ query: string }>
  ): Promise<void> {
    const query: string = event.detail.query;
    if (!query.trim()) return;

    searchQuery = query;
    isLoading = true;
    errorMessage = "";

    try {
        await refreshData(page, activeTab, searchQuery);
    } catch (error: unknown) {
      // console.error("Search error:", error);
      errorMessage = `Search error: ${(error as Error).message}`;
    } finally {
      isLoading = false;
    }
  }

  function handleVoiceError(event: CustomEvent<{ message: string }>): void {
    errorMessage = event.detail.message;
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
            newThreads = await fetchLatestThread(page+1, user.id, searchQuery);
          }

          if(!newThreads){
            hasMore = false;
            isLoading = false;
            return;
          }
          latestThreads = [...latestThreads, ...newThreads];
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

  onMount(async() => {
    observeScroll();
    await refreshData(page, activeTab, searchQuery);
  });

</script>

<main class="content-area" bind:this={contentAreaElement}>
  <div class="search-container">
    <SearchBar
      placeholder="Search with text or voice..."
      bind:value={searchQuery}
      on:search={handleSearch}
      on:error={handleVoiceError}
      on:input={handleSearch}
    />
  </div>

  <div class="tabs">
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "top"}
      on:click={() => {
        activeTab = "top";
        // refreshData(currentPage, searchQuery, activeTab, selectedCategory);
      }}
    >
      Top
    </button>
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "latest"}
      on:click={() => {
        activeTab = "latest";
        refreshData(page, activeTab, searchQuery);
      }}
    >
      Latest
    </button>
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "people"}
      on:click={() => {
        activeTab = "people";
        // refreshData(currentPage, searchQuery, activeTab, selectedCategory);
      }}
    >
      People
    </button>
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "media"}
      on:click={() => {
        activeTab = "media";
        // refreshData(currentPage, searchQuery, activeTab, selectedCategory);
      }}
    >
      Media
    </button>
    <button
      type="button"
      class="tab-button"
      class:active={activeTab === "communities"}
      on:click={() => {
        activeTab = "communities";
        // refreshData(currentPage, searchQuery, activeTab, selectedCategory);
      }}
    >
      Communities
    </button>
  </div>

  {#if errorMessage}
    <div class="error-message">
      {errorMessage}
    </div>
  {:else}
    {#if activeTab === "latest"}
        <ThreadFeed threads={latestThreads} {user} />
    <!-- {:else if activeTab === "communities"} -->

    {/if}

    <div bind:this={sentinel} class="loading-sentinel">
        {#if isLoading}
            <div class="skeleton-loading">
            <!-- eslint-disable-next-line @typescript-eslint/no-unused-vars -->
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
        {:else if !hasMore && latestThreads.length > 0}
            <p>No more thread available!</p>
        {:else if latestThreads.length == 0}
            <p>No thread available!</p>
        {/if}
    </div>
  {/if}
</main>


<style>
    .content-area {
      flex: 1;
      border-right: 1px solid var(--border-color);
      width: 50%;
      height: 100vh;
      margin: 0 auto;
      overflow-y: auto;
    }

    .search-container {
        margin: 2rem;
    }

    .tabs {
        display: flex;
        border-bottom: 1px solid var(--border-color);
    }

    .tab-button {
        flex: 1;
        background: none;
        border: none;
        color: var(--text-tertiary);
        font-size: 15px;
        font-weight: 500;
        padding: 16px;
        cursor: pointer;
        position: relative;
        transition: color 0.2s;
    }

    .tab-button.active {
        color: var(--button-primary);
        font-weight: 700;
    }

    .tab-button.active::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        height: 4px;
        background-color: var(--button-primary);
        border-radius: 4px 4px 0 0;
    }

    .error-message {
        padding: 1rem;
        background-color: #ffebee;
        color: #c62828;
        border-radius: 4px;
        margin-bottom: 1rem;
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

    @keyframes spin {
        0% {
            transform: rotate(0deg);
        }
        100% {
            transform: rotate(360deg);
        }
    }
</style>
