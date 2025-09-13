import React from 'react';

export interface ProblemItemProps {
  name: string;
  url: string;
  time: string;
  solved: boolean;
  onToggle: () => void;
}

const ProblemItem: React.FC<ProblemItemProps> = ({ name, url, time, solved, onToggle }) => (
  <li className="problem-item">
    <input type="checkbox" checked={solved} onChange={onToggle} />
    <span className={solved ? 'solved' : ''}>{name}</span>
    <a href={url} target="_blank" rel="noopener noreferrer" className="url-link">link</a>
    <span style={{ marginLeft: 8, color: '#888' }}>({time})</span>
  </li>
);

export default ProblemItem;
