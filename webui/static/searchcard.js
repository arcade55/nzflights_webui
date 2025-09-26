const searchTemplate = document.createElement('template');
searchTemplate.innerHTML = `
  <style>
@import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@24,400,0,0');
    .material-symbols-outlined {
      font-family: 'Material Symbols Outlined';
      font-weight: normal;
      font-style: normal;
      font-size: 24px; /* Adjust as needed */
      line-height: 1;
      letter-spacing: normal;
      text-transform: none;
      display: inline-block;
      white-space: nowrap;
      word-wrap: normal;
      direction: ltr;
      -webkit-font-smoothing: antialiased;
      text-rendering: optimizeLegibility;
      -moz-osx-font-smoothing: grayscale;
      font-feature-settings: 'liga';
    }

    :host {
      display: block; /* Ensures the component takes up block-level space */
      width: 100%; /* Fills the container, achieving the same width as other cards */
    }
    /* All styles are scoped to this component */
    .search-ticket {
      border-radius: 28px;
      padding: 24px;
      position: relative;
      background-color: var(--card-bg-focused, #D2E5E1);
      color: var(--card-text-focused, #00201B);
      font-family: 'Roboto', sans-serif;
    }
    /* Ticket stub cutout effect */
    .search-ticket::before,
    .search-ticket::after {
      content: '';
      position: absolute;
      width: 24px;
      height: 24px;
      border-radius: 50%;
      left: 50%;
      transform: translateX(-50%);
      /* It's crucial that the page provides this background color variable */
      background-color: var(--page-bg, #1C1C1E);
    }
    .search-ticket::before {
      top: -12px;
    }
    .search-ticket::after {
      bottom: -12px;
    }
    .input-field {
      background-color: var(--card-input-bg-focused, #FFFFFF);
      border-radius: 12px;
      padding: 12px 16px;
      margin-bottom: 12px;
      display: flex;
      align-items: center;
      gap: 12px;
    }
    .input-field input {
        border: none; outline: none; width: 100%;
        font-size: 16px; background: transparent;
        color: var(--card-text-focused, #00201B);
    }
    .input-field input::placeholder { color: var(--card-text-secondary-focused, #3F4946); }
    .input-field .icon {
      color: var(--card-text-secondary-focused, #3F4946);
    }
    .input-field .placeholder {
      font-size: 16px;
      color: var(--card-text-secondary-focused, #3F4946);
    }
    .input-field .label {
      font-size: 12px;
      color: var(--card-text-secondary-focused, #3F4946);
    }
    .input-field .value {
      font-size: 16px;
      font-weight: 500;
    }
    .separator {
      display: flex;
      align-items: center;
      text-align: center;
      color: var(--card-text-secondary-focused, #3F4946);
      margin: 24px 0;
      font-size: 12px;
      font-weight: 500;
      text-transform: uppercase;
    }
    .separator::before,
    .separator::after {
      content: '';
      flex: 1;
      border-bottom: 1px solid var(--card-outline-focused, #BFC9C5);
    }
    .separator:not(:empty)::before {
      margin-right: .75em;
    }
    .separator:not(:empty)::after {
      margin-left: .75em;
    }
    .action-button {
      width: 100%;
      background-color: var(--card-text-focused, #00201B);
      color: var(--card-bg-focused, #D2E5E1);
      border: none;
      padding: 16px;
      border-radius: 100px;
      font-size: 16px;
      font-weight: 500;
      cursor: pointer;
      display: flex;
      justify-content: center;
      align-items: center;
      gap: 8px;
      margin-top: 24px;
    }
  </style>
  <div class="search-ticket">
    <div class="input-field">
      <span class="material-symbols-outlined icon">airplane_ticket</span>


      <input type="text" id="flight-search-input" placeholder="Search by Flight Number...">





   

  </div>
`;
class SearchCard extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        this.shadowRoot.appendChild(searchTemplate.content.cloneNode(true));
        this.inputElement = this.shadowRoot.getElementById('flight-search-input');
    }
    connectedCallback() {
        // Listen for input changes and dispatch custom event
        // ... (rest of the code unchanged)
        this.inputElement.addEventListener('input', (e) => {
            const value = e.target.value;
            console.log(`Internal input event fired. Value: "${value}"`);
            // Dispatch custom event for Datastar to catch (lowercase)
            const event = new CustomEvent('searchtermchange', { // <-- Changed to lowercase
                bubbles: true,
                composed: true,
                detail: { value }
            });
            this.dispatchEvent(event);
        });
    }
}
customElements.define('search-card', SearchCard);