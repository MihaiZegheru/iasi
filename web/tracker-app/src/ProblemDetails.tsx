import React, { useEffect, useState } from 'react';
import AccordionBox from './AccordionBox';
import MarkdownView from './MarkdownView';
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
    <div className="problem-details-card">
      <div style={{ width: '100%' }}>
        <h2 className="problem-details-title">{problem.name}</h2>
      </div>
      <div className="problem-details-links-alt">
        <a
          href={problem.url}
          target="_blank"
          rel="noopener noreferrer"
          className="infoarena-btn"
        >
          <span className="icon" role="img" aria-label="Infoarena">🌐</span> Infoarena
        </a>
        <a
          href={`https://www.infoarena.ro/job_detail/${problem.id}?action=view-source`}
          target="_blank"
          rel="noopener noreferrer"
          className="infoarena-btn"
        >
          <span className="icon" role="img" aria-label="Code">📝</span> Code
        </a>
      </div>
      {editorial ? (
        <>
          <div className="problem-details-tabs" style={{position: 'relative'}}>
            <span
              className="flip-highlight"
              style={{
                transform: tab === 'editorial' ? 'translateX(100%)' : 'translateX(0%)',
              }}
            />
            <button onClick={() => setTab('hints')} disabled={tab === 'hints'}>Hints</button>
            <button onClick={() => setTab('editorial')} disabled={tab === 'editorial'}>Editorial</button>
          </div>
          {tab === 'hints' ? (
            <div className="problem-details-accordion">
              {editorial.hints.map((hint, i) => (
                <AccordionBox key={i} title={`Hint ${i + 1}`}>
                  <MarkdownView>{hint}</MarkdownView>
                </AccordionBox>
              ))}
            </div>
          ) : (
            <div className="problem-details-accordion">
              <AccordionBox title="Editorial">
                <MarkdownView>{editorial.editorial}</MarkdownView>
              </AccordionBox>
            </div>
          )}
        </>
      ) : (
        <>
          <button className="problem-details-generate-btn" onClick={handleGenerate} disabled={loading}>
            {loading ? 'Generating...' : 'Generate Hints/Editorial'}
          </button>
          {error && <div style={{ color: 'red', marginTop: 8 }}>{error}</div>}
        </>
      )}
      <div className="problem-details-back">
        <Link to="/">Back to list</Link>
      </div>
    </div>
  );
};

export default ProblemDetails;
