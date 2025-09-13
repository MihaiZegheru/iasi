import React from 'react';
import ProblemItem from './ProblemItem';
import type { ProblemItemProps } from './ProblemItem';

interface ProblemListProps {
  problems: ProblemItemProps[];
  solved: Record<string, boolean>;
  onToggle: (name: string) => void;
  filter: string;
}

const ProblemList: React.FC<ProblemListProps> = ({ problems, solved, onToggle, filter }) => {
  return (
    <ul className="problem-list">
      {problems
        .filter(p => p.name.toLowerCase().includes(filter.toLowerCase()))
        .map(p => (
          <ProblemItem
            key={p.name}
            {...p}
            solved={!!solved[p.name]}
            onToggle={() => onToggle(p.name)}
          />
        ))}
    </ul>
  );
};

export default ProblemList;
