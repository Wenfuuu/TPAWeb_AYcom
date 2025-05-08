<script lang="ts">
  // import { onMount } from "svelte";
  import SearchBar from "../lib/components/SearchBar.svelte";

  // State
  let searchQuery: string = "";
  let searchResults: { title: string; description: string }[] = [];
  let isLoading: boolean = false;
  let errorMessage: string = "";

  // Handle search
  async function handleSearch(
    event: CustomEvent<{ query: string }>
  ): Promise<void> {
    const query: string = event.detail.query;
    if (!query.trim()) return;

    searchQuery = query;
    isLoading = true;
    errorMessage = "";

    try {
      const response: Response = await fetch(
        `/api/search?q=${encodeURIComponent(searchQuery)}`
      );

      if (!response.ok) {
        throw new Error(`Search failed with status: ${response.status}`);
      }

      const data: { results: { title: string; description: string }[] } =
        await response.json();
      searchResults = data.results || [];
    } catch (error: unknown) {
      console.error("Search error:", error);
      errorMessage = `Search error: ${(error as Error).message}`;
      searchResults = [];
    } finally {
      isLoading = false;
    }
  }

  // Handle voice input error
  function handleVoiceError(event: CustomEvent<{ message: string }>): void {
    errorMessage = event.detail.message;
  }
</script>

<div class="search-page">
  <header class="search-header">
    <h1>Search</h1>
    <p>Use text or voice to search</p>
  </header>

  <div class="search-container">
    <SearchBar
      placeholder="Search with text or voice..."
      bind:value={searchQuery}
      on:search={handleSearch}
      on:error={handleVoiceError}
      autofocus={true}
    />
  </div>

  {#if errorMessage}
    <div class="error-message">
      {errorMessage}
    </div>
  {/if}

  {#if isLoading}
    <div class="loading">
      <span class="loader"></span>
      <p>Searching...</p>
    </div>
  {:else if searchResults.length > 0}
    <div class="search-results">
      <h2>Results for "{searchQuery}"</h2>
      <ul>
        {#each searchResults as result (result.description)}
          <li class="result-item">
            <h3>{result.title}</h3>
            <p>{result.description}</p>
          </li>
        {/each}
      </ul>
    </div>
  {:else if searchQuery}
    <div class="no-results">
      <p>No results found for "{searchQuery}"</p>
    </div>
  {/if}
</div>

<style>
  .search-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 2rem 1rem;
  }

  .search-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .search-header h1 {
    font-size: 2rem;
    margin-bottom: 0.5rem;
  }

  .search-container {
    margin-bottom: 2rem;
  }

  .error-message {
    padding: 1rem;
    background-color: #ffebee;
    color: #c62828;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin: 2rem 0;
  }

  .loader {
    width: 40px;
    height: 40px;
    border: 4px solid #f3f3f3;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 1rem;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .search-results h2 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
  }

  .result-item {
    padding: 1rem;
    border-bottom: 1px solid #eee;
    margin-bottom: 1rem;
  }

  .result-item h3 {
    font-size: 1.2rem;
    margin-bottom: 0.5rem;
  }

  .no-results {
    text-align: center;
    padding: 2rem;
    color: #666;
  }
</style>
