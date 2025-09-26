const template = document.createElement('template');
template.innerHTML = `
  <style>
    /* All styles for the flight card are now encapsulated within the component's Shadow DOM.
      They won't leak out or be affected by external styles.
    */
    @import url('https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@24,400,0,0');
    
    :host {
      display: block; /* By default, custom elements are inline */
      font-family: var(--font-family, 'Roboto', sans-serif);
    }

    .flight-card {
        background-color: var(--card-bg, #D2E5E1);
        border-radius: 28px;
        padding: 20px 24px;
        display: flex;
        flex-direction: column;
        gap: 16px;
        position: relative;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
        color: var(--card-text-primary, #00201B);
    }

    .flight-card::before,
    .flight-card::after {
        content: '';
        position: absolute;
        width: 24px;
        height: 24px;
        border-radius: 50%;
        left: 50%;
        transform: translateX(-50%);
        /* CSS variables from the main page can pierce the shadow DOM boundary */
        background-color: var(--background-color, #1C1C1E);
    }

    .flight-card::before {
        top: -12px;
    }

    .flight-card::after {
        bottom: -12px;
    }

    .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 12px;
    }

    .airline-info {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .airline-logo {
        width: 40px;
        height: 40px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;
        font-weight: 700;
        font-size: 16px;
    }

    /* Airline Brand Colors */
    .airline-nz { background-color: #000000; }
    .airline-mu { background-color: #1E3A8A; }
    .airline-qf { background-color: #E40000; }
    .airline-ek { background-color: #d81921; }
    .airline-sq { background-color: #F99F00; }
    .airline-fj { background-color: #2F2E2E; }
    .airline-cx { background-color: #006442; }
    .airline-ua { background-color: #002244; }
    .airline-ac { background-color: #F01428; }
    .airline-ke { background-color: #0064A2; }
    .airline-jq { background-color: #FF5500; }

    .flight-number {
        font-weight: 500;
        color: var(--card-text-primary, #00201B);
        font-size: 28px;
    }

    .airline-name {
        font-size: 14px;
        color: var(--card-text-secondary, #3F4946);
    }

    .card-share-button {
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
        color: var(--card-text-secondary, #3F4946);
    }

    .card-share-button .material-symbols-outlined {
        font-size: 24px;
    }

    .card-separator {
        border: 0;
        border-top: 1px dashed var(--card-separator-color, rgba(63, 73, 70, 0.4));
        margin: -4px 0;
    }

    .flight-path {
        display: grid;
        grid-template-columns: 1fr auto 1fr;
        align-items: center;
    }

    .location h2 {
        font-size: 28px;
        font-weight: 500;
        margin: 0;
        color: var(--card-text-primary, #00201B);
    }

    .location p {
        font-size: 14px;
        margin: 2px 0 0 0;
        color: var(--card-text-secondary, #3F4946);
    }

    .path-icon {
        color: var(--card-text-secondary, #3F4946);
    }

    .flight-details {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
    }

    .flight-path > .location:first-child,
    .flight-details > .detail-item:first-child {
        text-align: left;
    }

    .flight-path > .location:last-child,
    .flight-details > .detail-item:last-child {
        text-align: right;
    }

    .detail-item h3 {
        margin: 0;
        font-size: 20px;
        font-weight: 500;
        color: var(--card-text-primary, #00201B);
    }

    .detail-item p {
        margin: 2px 0 0 0;
        font-size: 12px;
        font-weight: 500;
        color: var(--card-text-secondary, #3F4946);
        text-transform: uppercase;
    }

    .card-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 14px;
        color: var(--card-text-secondary, #3F4946);
    }

    .status {
        font-weight: 500;
    }
    
    .status-ontime {
        color: var(--card-status-ontime, #1E8449);
    }

    .status-delayed {
        color: var(--card-status-delayed, #C0392B);
    }
  </style>

  <div class="flight-card">
      <div class="card-header">
          <div class="airline-info">
              <div id="logo" class="airline-logo"></div>
              <div>
                  <div id="flight-number" class="flight-number"></div>
                  <div id="airline-name" class="airline-name"></div>
              </div>
          </div>
          <button class="card-share-button">
              <span class="material-symbols-outlined">share</span>
          </button>
      </div>
      <hr class="card-separator">
      <div class="flight-path">
          <div class="location">
              <h2 id="origin-iata"></h2>
              <p id="origin-city"></p>
          </div>
          <div class="path-icon"><span class="material-symbols-outlined">east</span></div>
          <div class="location">
              <h2 id="dest-iata"></h2>
              <p id="dest-city"></p>
          </div>
      </div>
      <div class="flight-details">
          <div class="detail-item">
              <h3 id="gate"></h3>
              <p>Gate</p>
          </div>
          <div class="detail-item">
              <h3 id="boarding"></h3>
              <p>Boarding</p>
          </div>
      </div>
      <hr class="card-separator">
      <div class="card-footer">
          <span id="dep-time"></span>
          <span id="status" class="status"></span>
          <span id="arr-time"></span>
      </div>
  </div>
`;

class FlightCard extends HTMLElement {
  constructor() {
    super();
    // Attach a shadow root to the element.
    this.attachShadow({ mode: 'open' });
    this.shadowRoot.appendChild(template.content.cloneNode(true));
  }

  // This method is called when the element is added to the DOM.
  connectedCallback() {
    this._updateRendering();
  }

  // This method is called when an observed attribute changes.
  attributeChangedCallback(name, oldValue, newValue) {
      if (oldValue !== newValue) {
          this._updateRendering();
      }
  }

  // Define which attributes to watch for changes.
  static get observedAttributes() {
      return [
          'airline-logo-text', 'airline-class', 'flight-number', 'airline-name',
          'origin-iata', 'origin-city', 'dest-iata', 'dest-city', 'gate',
          'boarding-time', 'departure-time', 'status-text', 'status-class', 'arrival-time'
      ];
  }

  // Helper function to read attributes and update the shadow DOM.
  _updateRendering() {
    const setContent = (id, value) => {
        const element = this.shadowRoot.getElementById(id);
        if (element) {
            element.textContent = value || '';
        }
    };
    
    // Populate simple text content
    setContent('flight-number', this.getAttribute('flight-number'));
    setContent('airline-name', this.getAttribute('airline-name'));
    setContent('origin-iata', this.getAttribute('origin-iata'));
    setContent('origin-city', this.getAttribute('origin-city'));
    setContent('dest-iata', this.getAttribute('dest-iata'));
    setContent('dest-city', this.getAttribute('dest-city'));
    setContent('gate', this.getAttribute('gate'));
    setContent('boarding', this.getAttribute('boarding-time'));
    setContent('dep-time', this.getAttribute('departure-time'));
    setContent('arr-time', this.getAttribute('arrival-time'));
    setContent('status', this.getAttribute('status'));


    // Handle elements with classes that need to be set dynamically
    const logo = this.shadowRoot.getElementById('logo');
    if (logo) {
        logo.textContent = this.getAttribute('airline-logo-text') || '';
        logo.className = 'airline-logo'; // Reset classes
        const airlineClass = this.getAttribute('airline-class');
        if (airlineClass) {
            logo.classList.add(airlineClass);
        }
    }

    const status = this.shadowRoot.getElementById('status');
    if (status) {
        status.textContent = this.getAttribute('status-text') || '';
        status.className = 'status'; // Reset classes
        const statusClass = this.getAttribute('status-class');
        if (statusClass) {
            status.classList.add(statusClass);
        }
    }
  }
}

// Define the new custom element so the browser knows what to do with the <flight-card> tag.
customElements.define('flight-card', FlightCard);

