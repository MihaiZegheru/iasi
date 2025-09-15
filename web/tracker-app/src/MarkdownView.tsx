import React from 'react';
import ReactMarkdown from 'react-markdown';

interface MarkdownViewProps {
  children: string;
}

const MarkdownView: React.FC<MarkdownViewProps> = ({ children }) => {
  return (
    <div className="markdown-centered">
      <ReactMarkdown>{children}</ReactMarkdown>
    </div>
  );
};

export default MarkdownView;
