import React, { useEffect, useState } from 'react';
import ProblemList from './ProblemList';
import type { Problem } from './types';
import './App.css';

interface ProblemsResponse {
  username: string;
  problems: Problem[];
}

const LOCAL_KEY_PREFIX = 'iasi_tracker_';

const App: React.FC = () => {
  const [problems, setProblems] = useState<Problem[]>([]);
  const [username, setUsername] = useState('');
  const [solved, setSolved] = useState<Record<string, boolean>>({});
  const [filter, setFilter] = useState('');

  useEffect(() => {
    fetch('/problems')
      .then(r => r.json())
      .then((data: ProblemsResponse) => {
        setProblems(data.problems);
        setUsername(data.username);
        const saved = localStorage.getItem(LOCAL_KEY_PREFIX + data.username);
        if (saved) setSolved(JSON.parse(saved));
      });
  }, []);

  const handleToggle = (name: string) => {
    setSolved(prev => {
      const next = { ...prev, [name]: !prev[name] };
      localStorage.setItem(LOCAL_KEY_PREFIX + username, JSON.stringify(next));
      return next;
    });
  };

  const solvedCount = problems.filter(p => solved[p.name]).length;

  return (
    <div className="tracker-container">
  <h1>Infoarena Scout &amp; Index</h1>
      <div className="progress">Solved: {solvedCount} / {problems.length}</div>
      <input
        type="text"
        placeholder="Search problems..."
        value={filter}
        onChange={e => setFilter(e.target.value)}
        style={{ marginBottom: '1em', width: '100%', padding: '0.5em' }}
      />
      <ProblemList
        problems={problems}
        solved={solved}
        onToggle={handleToggle}
        filter={filter}
      />
    </div>
  );
};

export default App;
