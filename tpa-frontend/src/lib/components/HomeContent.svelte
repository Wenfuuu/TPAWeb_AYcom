<script lang="ts">
  import { onMount } from "svelte";
  import { type User, type UserThread } from "../model/user";
  import {
    fetchFollowingThread,
    fetchThreadWithUser,
    type Thread,
  } from "../model/thread";

  import ComposePost from "./ComposePost.svelte";
  import ThreadFeed from "./ThreadFeed.svelte";

  export let threads: (Thread & { userData: UserThread })[] = [];

  export let refreshData: (page: number) => Promise<void>;
  export let refreshFollowingData: (page: number) => Promise<void>;

  export let user: User | null;

  let page = 1;
  let isLoading = false;
  let hasMore = true;

  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;

  export let currentActiveTab: "for-you" | "following";

  function switchTab(tab: "for-you" | "following") {
    if (currentActiveTab !== tab) {
      currentActiveTab = tab;
      page = 1;
      threads = [];
      hasMore = true;

      console.log(currentActiveTab);
      if(currentActiveTab === "for-you") refreshData(page);
      else refreshFollowingData(page);
    }
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
        console.log(hasMore);
        if (entries[0].isIntersecting && !isLoading && hasMore) {
          isLoading = true;
          console.log("Loading more content, page:", page + 1);

          await new Promise(resolve => setTimeout(resolve, 500));

          let newThreads;
          if(user){
            if(currentActiveTab === "for-you") {
              newThreads = await fetchThreadWithUser(page+1, user.id);
            }else{
              newThreads = await fetchFollowingThread(page+1, user.id);
            }
          }

          if(!newThreads){
            hasMore = false;
            isLoading = false;
            return;
          }
          console.log("Fetched new threads:", newThreads.length);
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
      console.log(
        "Observer attached to sentinel element with content area root."
      );
    } else {
      console.warn(
        "Sentinel element not found - scroll observation not active"
      );
    }
  }

  async function handlePostCreated() {
    hasMore = true;
    if(currentActiveTab === "for-you") await refreshData(1);
    else refreshFollowingData(1);
  }

  onMount(async () => {
    observeScroll();
    await refreshData(page);
  });
</script>

{#if user}
  <main class="content-area" bind:this={contentAreaElement}>
    <div class="header">
      <div class="tabs">
        <button
          type="button"
          class="tab"
          class:active={currentActiveTab === "for-you"}
          on:click={() => switchTab("for-you")}
        >
          For you
        </button>
        <button
          type="button"
          class="tab"
          class:active={currentActiveTab === "following"}
          on:click={() => switchTab("following")}
        >
          Following
        </button>
      </div>
    </div>

    <ComposePost {user} refreshData={handlePostCreated} />

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
      {/if}
    </div>
  </main>
{/if}

<style>
  /* Styles for HomeContent - Keep only styles relevant to this component's structure */
  .content-area {
    flex: 1;
    border-right: 1px solid var(--border-color);
    width: 50%;
    height: 100vh;
    margin: 0 auto; /* Center if using max-width */
    overflow-y: auto;
  }

  .header {
    position: sticky;
    top: 0;
    background-color: var(--card-background);
    backdrop-filter: blur(12px);
    z-index: 10;
    border-bottom: 1px solid var(--border-color);
  }

  .tabs {
    display: flex;
    height: 53px;
  }

  .tab {
    flex: 1;
    background: none;
    border: none;
    color: var(--text-secondary); /* Default color */
    font-size: 15px;
    font-weight: 500;
    padding: 16px;
    cursor: pointer;
    position: relative;
    transition:
      background-color 0.2s,
      color 0.2s;
  }

  .tab:hover {
    background-color: var(--hover-background);
    color: var(--text-primary);
  }

  .tab.active {
    color: var(--text-primary);
    font-weight: 700;
  }

  .tab.active::after {
    content: "";
    position: absolute;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 56px;
    height: 4px;
    background-color: var(--button-primary);
    border-radius: 9999px;
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

  @media (max-width: 768px) {
    .content-area {
      border-right: none; /* Remove border on mobile */
    }
  }
</style>
