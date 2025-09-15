import React, { useEffect, useState } from 'react';
import { Routes, Route } from 'react-router-dom';
import ProblemList from './ProblemList';
import ProblemDetails from './ProblemDetails';
import type { Problem } from './types';
import './App.css';

interface ProblemsResponse {
  username: string;
  problems: Problem[];
}

const GLOBAL_SOLVED_KEY = 'iasi_tracker_global_solved';

const App: React.FC = () => {
  const [problems, setProblems] = useState<Problem[]>([]);
  const [username, setUsername] = useState('');
  const [solved, setSolved] = useState<Record<string, boolean>>({});
  const [filter, setFilter] = useState('');
  const [sortOption, setSortOption] = useState('time-asc');

  useEffect(() => {
    fetch('/problems')
      .then(r => r.json())
      .then((data: ProblemsResponse) => {
        setProblems(data.problems);
        setUsername(data.username);
        const saved = localStorage.getItem(GLOBAL_SOLVED_KEY);
        if (saved) setSolved(JSON.parse(saved));
      });
  }, []);

  const handleToggle = (name: string) => {
    setSolved(prev => {
      const next = { ...prev, [name]: !prev[name] };
      localStorage.setItem(GLOBAL_SOLVED_KEY, JSON.stringify(next));
      return next;
    });
  };


  // Sorting logic
  const sortedProblems = React.useMemo(() => {
    let arr = [...problems];
    switch (sortOption) {
      case 'solved':
        arr.sort((a, b) => Number(!solved[a.name]) - Number(!solved[b.name]));
        break;
      case 'unsolved':
        arr.sort((a, b) => Number(!solved[b.name]) - Number(!solved[a.name]));
        break;
      case 'az':
        arr.sort((a, b) => a.name.localeCompare(b.name));
        break;
      case 'za':
        arr.sort((a, b) => b.name.localeCompare(a.name));
        break;
      case 'time-desc':
        arr.sort((a, b) => new Date(b.time).getTime() - new Date(a.time).getTime());
        break;
      case 'time-asc':
      default:
        arr.sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
        break;
    }
    return arr;
  }, [problems, solved, sortOption]);

  const solvedCount = problems.filter(p => solved[p.name]).length;

  return (
    <Routes>
      <Route
        path="/"
        element={
          <div className="tracker-container">
            <h1 className="gradient-title" style={{marginBottom: 0}}>Infoarena</h1>
            <h1 className="gradient-title" style={{fontSize: '2.2em', marginTop: 0.1 + 'em'}}>Scout &amp; Index</h1>
            <div className="progress">Solved: {solvedCount} / {problems.length}</div>
            <div style={{ display: 'flex', gap: 12, marginBottom: '1em' }}>
              <input
                type="text"
                placeholder="Search problems..."
                value={filter}
                onChange={e => setFilter(e.target.value)}
                className="tracker-input"
              />
              <select value={sortOption} onChange={e => setSortOption(e.target.value)} className="sort-dropdown">
                <option value="time-asc">Sort: Time ↑ (default)</option>
                <option value="time-desc">Time ↓</option>
                <option value="solved">Solved first</option>
                <option value="unsolved">Unsolved first</option>
                <option value="az">A-Z</option>
                <option value="za">Z-A</option>
              </select>
            </div>
            <ProblemList
              problems={sortedProblems}
              solved={solved}
              onToggle={handleToggle}
              filter={filter}
            />
          </div>
        }
      />
      <Route path="/problem/:id" element={<ProblemDetails />} />
    </Routes>
  );
};

export default App;
