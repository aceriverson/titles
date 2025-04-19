import styles from './App.module.css';

function App() {
  return (
    <div class={styles.App}>
      <header class={styles.header}>
        <div>
          <p>
            <code>titles.run</code> is undergoing some changes
          </p>
          <p>
            Please check back later
          </p>
        </div>
        <a
          class={styles.link}
          href="./polygons"
        >
          Legacy Site
        </a>
      </header>
    </div>
  );
}

export default App;
