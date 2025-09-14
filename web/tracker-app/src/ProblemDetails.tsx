import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Problem } from './types';

interface EditorialData {
  hints: string[];
  editorial: string;
}

const ProblemDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [problem, setProblem] = useState<Problem | null>(null);
  const [editorial, setEditorial] = useState<EditorialData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tab, setTab] = useState<'hints' | 'editorial'>('hints');

  useEffect(() => {
    fetch(`/problems`)
      .then(r => r.json())
      .then(data => {
        const found = data.problems.find((p: Problem) => p.id === id);
        setProblem(found || null);
      });
    fetch(`/problems/${id}/editorial`)
      .then(r => {
        if (!r.ok) throw new Error('Not generated');
        return r.json();
      })
      .then(setEditorial)
      .catch(() => setEditorial(null));
  }, [id]);

  const handleGenerate = async () => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`/problems/${id}/generate`, { method: 'POST' });
      if (!res.ok) throw new Error('Failed to generate');
      const data = await res.json();
      setEditorial(data);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  if (!problem) return <div>Problem not found. <Link to="/">Back</Link></div>;

  return (
    <div style={{ maxWidth: 700, margin: '0 auto' }}>
      <h2>{problem.name}</h2>
      <div style={{ marginBottom: 16 }}>
        <a href={problem.url} target="_blank" rel="noopener noreferrer">View on Infoarena</a>
      </div>
      {editorial ? (
        <>
          <div style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
            <button onClick={() => setTab('hints')} disabled={tab === 'hints'}>Hints</button>
            <button onClick={() => setTab('editorial')} disabled={tab === 'editorial'}>Editorial</button>
          </div>
          {tab === 'hints' ? (
            <ul>
              {editorial.hints.map((hint, i) => <li key={i}>{hint}</li>)}
            </ul>
          ) : (
            <div style={{ whiteSpace: 'pre-line', background: '#f8f8f8', padding: 16, borderRadius: 8 }}>{editorial.editorial}</div>
          )}
        </>
      ) : (
        <>
          <button onClick={handleGenerate} disabled={loading}>
            {loading ? 'Generating...' : 'Generate Hints/Editorial'}
          </button>
          {error && <div style={{ color: 'red', marginTop: 8 }}>{error}</div>}
        </>
      )}
      <div style={{ marginTop: 24 }}>
        <Link to="/">Back to list</Link>
      </div>
    </div>
  );
};

export default ProblemDetails;
