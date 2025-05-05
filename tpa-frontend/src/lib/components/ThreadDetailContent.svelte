<script lang="ts">
  import { onMount } from "svelte";
  import {
    bookmarkThread,
    fetchThreadRepliesWithUser,
    likeThread,
    unbookmarkThread,
    unlikeThread,
    type Thread,
  } from "../model/thread";
  import type { User, UserThread } from "../model/user";
  import { navigate } from "svelte-routing";
  import ComposeReply from "./ComposeReply.svelte";

  export let thread: Thread & { userData: UserThread };
  export let replies: (Thread & { userData: UserThread })[] = [];

  export let refreshData: (page: number) => Promise<void>;
  export let loadThreadData: (threadId: string, page: number) => Promise<void>;
  export let handleBackNavigation: () => void;
//   export let changePrevId: () => void;

  export let user: User | null;

  let page = 1;
  let isLoading = false;
  let hasMore = true;

  let contentAreaElement: HTMLElement | null = null;

  let sentinel: HTMLDivElement | null = null;
  let observer: IntersectionObserver | null = null;

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

          await new Promise((resolve) => setTimeout(resolve, 500));

          let newReplies;
          if (user) {
            newReplies = await fetchThreadRepliesWithUser(
              page+1,
              thread.id,
              user.id
            );
          }

          if (!newReplies) {
            hasMore = false;
            isLoading = false;
            return;
          }
          console.log("Fetched new replies:", newReplies.length);
          replies = [...replies, ...newReplies];
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

  async function handleReplyCreated() {
    hasMore = true;
    loadThreadData(thread.id, page);
  }

  function toggleLike(targetThread: Thread) {
    targetThread.isLiked = !targetThread.isLiked;
  }

  async function handleLikeThread(): Promise<void> {
    if (!user) return;

    if (!thread.isLiked) {
      await likeThread(user.id, thread.id);
    } else {
      await unlikeThread(user.id, thread.id);
    }

    toggleLike(thread);
    thread = { ...thread }; // Force UI update
  }

  async function handleLikeReplies(replyThread: Thread): Promise<void> {
    if (!user) return;

    if (!replyThread.isLiked) {
      await likeThread(user.id, replyThread.id);
    } else {
      await unlikeThread(user.id, replyThread.id);
    }

    toggleLike(replyThread);
    replies = [...replies]; // Force UI update
  }

  function toggleBookmark(targetThread: Thread) {
    targetThread.isBookmarked = !targetThread.isBookmarked;
  }

  async function handleBookmarkThread(): Promise<void> {
    if (!user) return;

    if (!thread.isBookmarked) {
      await bookmarkThread(user.id, thread.id);
    } else {
      await unbookmarkThread(user.id, thread.id);
    }

    toggleBookmark(thread);
    thread = { ...thread }; // Force UI update
  }

  async function handleBookmarkReplies(replyThread: Thread): Promise<void> {
    if (!user) return;

    if (!replyThread.isBookmarked) {
      await bookmarkThread(user.id, replyThread.id);
    } else {
      await unbookmarkThread(user.id, replyThread.id);
    }

    toggleBookmark(replyThread);
    replies = [...replies]; // Force UI update
  }

  onMount(async () => {
    observeScroll();
    if(thread) loadThreadData(thread.id, 1);
  });
</script>

{#if user && thread}
  <main class="content-area" bind:this={contentAreaElement}>
    <!-- Header -->
    <div class="header">
      <button type="button" class="back-button" aria-label="Back button" on:click={() => {
        handleBackNavigation();
      }}>
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="m15 18-6-6 6-6"></path>
        </svg>
      </button>
      <h1>Post</h1>
    </div>

    <!-- Main Thread -->
    <div class="thread-post">
      <div class="post-avatar">
        <img src={thread.userData?.profile || "/placeholder.svg"} alt={thread.userData?.name || "User"} />
      </div>
      <div class="post-content">
        <div class="post-header">
          <div class="post-user-info">
            <div class="user-name-container">
              <span class="post-user-name">{thread.userData?.name || "Unknown User"}</span>
              {#if thread.userData?.isPremium}
                <span class="verified-badge">
                  <svg viewBox="0 0 24 24" width="16" height="16" fill="rgb(29, 155, 240)">
                    <path d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"></path>
                  </svg>
                </span>
              {/if}
            </div>
            <span class="post-username">@{thread.userData?.username || "unknown"}</span>
          </div>
          <button type="button" class="post-more" aria-label="Post More">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="1"></circle>
              <circle cx="19" cy="12" r="1"></circle>
              <circle cx="5" cy="12" r="1"></circle>
            </svg>
          </button>
        </div>

        {#if thread.content}
          <div class="post-text">{thread.content}</div>
        {/if}

        {#if thread.images && thread.images.length > 0}
          <div class="post-images" class:grid-2={thread.images.length === 2} class:grid-3={thread.images.length === 3} class:grid-4={thread.images.length === 4}>
            {#each thread.images as image, i (i)}
              <div class="image-container">
                <!-- <img src={image || "/placeholder.svg"} alt="Post image {i + 1}" /> -->
                 {#if image.endsWith(".png") || image.endsWith(".jpg")
                 || image.endsWith(".jpeg") || image.endsWith(".gif")}
                  <img
                    src={image || "/placeholder.svg"}
                    alt="Post image {i + 1}"
                  />
                {:else}
                  <video src={image} controls class="post-video">
                    <track kind="captions" src="path_to_captions.vtt" srclang="en" label="English">
                  </video>
                {/if}
              </div>
            {/each}
          </div>
        {/if}

        <div class="post-time">{thread.created_at}</div>

        <div class="post-actions">
          <button type="button" class="post-action reply" class:active={false}>
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
            </svg>
            <span>{thread.insights?.replies || 0}</span>
          </button>
          <button type="button" class="post-action repost">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="m17 2 4 4-4 4"></path>
              <path d="M3 11v-1a4 4 0 0 1 4-4h14"></path>
              <path d="m7 22-4-4 4-4"></path>
              <path d="M21 13v1a4 4 0 0 1-4 4H3"></path>
            </svg>
            <span>{thread.insights?.reposts || 0}</span>
          </button>
          <button type="button" class="post-action like" class:active={thread.isLiked} on:click|stopPropagation={handleLikeThread}>
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"></path>
            </svg>
            <span>{thread.insights?.likes || 0}</span>
          </button>
          <button type="button" class="post-action bookmark" class:active={thread.isBookmarked} on:click|stopPropagation={handleBookmarkThread} aria-label="Bookmark Post">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Reply Composer -->
    <ComposeReply {thread} {user} refreshData={handleReplyCreated} />

    <div class="posts">
      {#each replies as replyThread (replyThread.id)}
        <div
          class="post"
          role="button"
          tabindex="0"
          on:click={() => {
            navigate(`/thread-detail/${replyThread.id}`);
            refreshData(1);
          }}
          on:keydown={(e) => {
            if (e.key === "Enter" || e.key === " ") {
              navigate(`/thread-detail/${replyThread.id}`);
            }
          }}
        >
          <div class="post-avatar">
            <img
              src={replyThread.userData?.profile || "/placeholder.svg"}
              alt={replyThread.userData?.name || "User"}
            />
          </div>
          <div class="post-content">
            <div class="reply-to">
              Replied to <span class="reply-username">@{thread.userData?.username || "unknown"}</span>
            </div>
            <div class="post-header">
              <div class="post-reply-info">
                <span class="post-user-name"
                  >{replyThread.userData?.name || "Unknown User"}</span
                >
                {#if replyThread.userData?.isPremium}
                  <span class="verified-badge">
                    <svg
                      viewBox="0 0 24 24"
                      width="16"
                      height="16"
                      fill="rgb(29, 155, 240)"
                    >
                      <path
                        d="M22.25 12c0-1.43-.88-2.67-2.19-3.34.46-1.39.2-2.9-.81-3.91s-2.52-1.27-3.91-.81c-.66-1.31-1.91-2.19-3.34-2.19s-2.67.88-3.33 2.19c-1.4-.46-2.91-.2-3.92.81s-1.26 2.52-.8 3.91c-1.31.67-2.2 1.91-2.2 3.34s.89 2.67 2.2 3.34c-.46 1.39-.21 2.9.8 3.91s2.52 1.26 3.91.81c.67 1.31 1.91 2.19 3.34 2.19s2.68-.88 3.34-2.19c1.39.45 2.9.2 3.91-.81s1.27-2.52.81-3.91c1.31-.67 2.19-1.91 2.19-3.34zm-11.71 4.2L6.8 12.46l1.41-1.42 2.26 2.26 4.8-5.23 1.47 1.36-6.2 6.77z"
                      ></path>
                    </svg>
                  </span>
                {/if}
                <span class="post-username"
                  >@{replyThread.userData?.username || "unknown"}</span
                >
                <span class="post-date">· {replyThread.created_at}</span>
              </div>
              <button type="button" class="post-more" aria-label="Post More">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <circle cx="12" cy="12" r="1"></circle>
                  <circle cx="19" cy="12" r="1"></circle>
                  <circle cx="5" cy="12" r="1"></circle>
                </svg>
              </button>
            </div>

            {#if replyThread.content}
              <div class="post-text">{replyThread.content}</div>
            {/if}

            {#if replyThread.images && replyThread.images.length > 0}
              <div
                class="post-images"
                class:grid-2={replyThread.images.length === 2}
                class:grid-3={replyThread.images.length === 3}
                class:grid-4={replyThread.images.length === 4}
              >
                {#each replyThread.images as image, i (i)}
                  <div class="image-container">
                    <img
                      src={image || "/placeholder.svg"}
                      alt="Post image {i + 1}"
                    />
                  </div>
                {/each}
              </div>
            {/if}

            <div class="post-actions">
              <button
                type="button"
                class="post-action reply"
                class:active={false}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
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
                <span>{replyThread.insights?.replies || 0}</span>
              </button>
              <button type="button" class="post-action repost">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="m17 2 4 4-4 4"></path>
                  <path d="M3 11v-1a4 4 0 0 1 4-4h14"></path>
                  <path d="m7 22-4-4 4-4"></path>
                  <path d="M21 13v1a4 4 0 0 1-4 4H3"></path>
                </svg>
                <span>{replyThread.insights?.reposts || 0}</span>
              </button>
              <button
                type="button"
                class="post-action like"
                class:active={replyThread.isLiked}
                on:click|stopPropagation={() => handleLikeReplies(replyThread)}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path
                    d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"
                  ></path>
                </svg>
                <span>{replyThread.insights?.likes || 0}</span>
              </button>
              <button
                type="button"
                class="post-action bookmark"
                class:active={replyThread.isBookmarked}
                on:click|stopPropagation={() => handleBookmarkReplies(replyThread)}
                aria-label="Bookmark Post"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="18"
                  height="18"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M19 21l-7-5-7 5V5a2 2 0 0 1 2-2h10a2 2 0 0 1 2 2v16z"
                  ></path>
                </svg>
              </button>
            </div>
          </div>
        </div>
      {/each}
    </div>

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
      {:else if !hasMore && replies.length > 0}
        <p>No more replies available!</p>
      {:else if replies.length === 0 && !isLoading}
        <div class="no-replies">
          <p>No replies yet. Be the first to reply!</p>
        </div>
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
    margin: 0 auto;
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
    display: flex;
    align-items: center;
    padding: 12px 16px;
    gap: 20px;
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

  .header-actions {
    margin-left: auto;
    display: flex;
    gap: 8px;
  }

  .reply-button {
    background-color: var(--button-primary);
    color: var(--button-primary-text);
    border: none;
    border-radius: 9999px;
    padding: 8px 16px;
    font-size: 15px;
    font-weight: 700;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .reply-button:hover {
    background-color: var(--button-primary-hover);
  }

  .options-button {
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

  .options-button:hover {
    background-color: var(--hover-background);
  }

  /* Main thread post styles */
  .thread-post {
    display: flex;
    padding: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .post-avatar {
    margin-right: 12px;
    flex-shrink: 0;
  }

  .thread-post .post-avatar img {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }

  .post-avatar img {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .post-content {
    flex: 1;
    min-width: 0;
  }

  .post-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 4px;
  }

  .post-user-info {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .thread-post .post-user-info {
    display: flex;
    flex-direction: column;
  }

  .user-name-container {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .post-user-name {
    font-weight: 700;
    color: var(--text-primary);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .thread-post .post-user-name {
    font-size: 16px;
  }

  .verified-badge {
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }

  .post-username,
  .post-date {
    color: var(--text-tertiary);
    white-space: nowrap;
    font-size: 14px;
  }

  .thread-post .post-username {
    font-size: 15px;
  }

  .post-username {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .post-more {
    background: none;
    border: none;
    color: var(--text-tertiary);
    cursor: pointer;
    padding: 4px;
    margin-left: 8px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
      background-color 0.2s,
      color 0.2s;
    width: 30px;
    height: 30px;
    flex-shrink: 0;
  }

  .post-more:hover {
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
    color: var(--button-primary);
  }

  .post-text {
    color: var(--text-primary);
    margin-bottom: 12px;
    word-wrap: break-word;
    line-height: 1.5;
    font-size: 15px;
  }

  .thread-post .post-text {
    font-size: 17px;
  }

  .post-time {
    color: var(--text-tertiary);
    font-size: 15px;
    margin: 16px 0;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--border-color);
  }

  .post-images {
    border-radius: 16px;
    overflow: hidden;
    margin-top: 8px;
    margin-bottom: 12px;
    display: grid;
    gap: 2px;
    border: 1px solid var(--border-color);
  }

  /* Grid layouts for different image counts */
  .post-images.grid-2 {
    grid-template-columns: repeat(2, 1fr);
    aspect-ratio: 16 / 9;
  }
  .post-images.grid-3 {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    aspect-ratio: 1 / 1;
  }
  .post-images.grid-3 .image-container:first-child {
    grid-column: 1 / 2;
    grid-row: 1 / 3;
  }
  .post-images.grid-3 .image-container:nth-child(2) {
    grid-column: 2 / 3;
    grid-row: 1 / 2;
  }
  .post-images.grid-3 .image-container:nth-child(3) {
    grid-column: 2 / 3;
    grid-row: 2 / 3;
  }
  .post-images.grid-4 {
    grid-template-columns: repeat(2, 1fr);
    grid-template-rows: repeat(2, 1fr);
    aspect-ratio: 1 / 1;
  }

  .image-container {
    position: relative;
    overflow: hidden;
    background-color: var(--hover-background);
  }

  .image-container img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: opacity 0.3s ease;
  }

  .post-actions {
    display: flex;
    justify-content: space-between;
    max-width: 425px;
    margin-top: 12px;
  }

  .post-action {
    background: none;
    border: none;
    color: var(--text-tertiary);
    display: flex;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    transition: color 0.2s;
    font-size: 13px;
    padding: 6px;
    border-radius: 4px;
  }

  .post-action.reply:hover {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
  }

  .post-action.reply.active {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.5);
  }

  .post-action.like:hover {
    color: red;
    background-color: rgba(var(--button-primary-rgb, 200, 15, 24), 0.1);
  }

  .post-action.like.active {
    color: red;
    background-color: rgba(var(--button-primary-rgb, 200, 15, 24), 0.7);
  }

  .post-action.bookmark:hover {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.1);
  }

  .post-action.bookmark.active {
    color: var(--button-primary);
    background-color: rgba(var(--button-primary-rgb, 29, 155, 240), 0.7);
  }

  .post-action svg {
    width: 18px;
    height: 18px;
  }

  /* Replies section */
  .reply-to {
    color: var(--text-tertiary);
    font-size: 14px;
    margin-bottom: 4px;
  }

  .reply-username {
    color: var(--button-primary);
  }

  .no-replies {
    padding: 32px 16px;
    text-align: center;
    color: var(--text-tertiary);
  }

  .posts {
    display: flex;
    flex-direction: column;
  }

  .post {
    display: flex;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-color);
    transition: background-color 0.2s;
    cursor: pointer;
  }

  .post:hover {
    background-color: var(--hover-background);
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

  .post-video {
    position: relative;
    width: 100%;
    height: 280px;
    background-color: #000;
    overflow: hidden;
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
      border-right: none;
      width: 100%;
    }
    .post-actions {
      max-width: 100%;
    }
  }

  @media (max-width: 500px) {
    .post-user-info {
      font-size: 14px;
    }
    .post-text {
      font-size: 14px;
    }
    .post-action span {
      display: none;
    }
    .post-action {
      gap: 4px;
    }
    .thread-post .post-text {
      font-size: 16px;
    }
  }
</style>