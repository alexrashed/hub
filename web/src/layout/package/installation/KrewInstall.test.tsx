import { render } from '@testing-library/react';
import React from 'react';

import { Repository } from '../../../types';
import KrewInstall from './KrewInstall';

const repo: Repository = {
  kind: 5,
  name: 'repo',
  displayName: 'Repo',
  url: 'http://repo.test',
  userAlias: 'user',
};
const defaultProps = {
  name: 'packageName',
  repository: repo,
};

describe('KrewInstall', () => {
  it('creates snapshot', () => {
    const result = render(<KrewInstall {...defaultProps} />);
    expect(result.asFragment()).toMatchSnapshot();
  });

  describe('Render', () => {
    it('renders component', () => {
      const { getByText } = render(<KrewInstall {...defaultProps} />);

      expect(getByText('Add repository')).toBeInTheDocument();
      expect(getByText('kubectl krew index add repo http://repo.test')).toBeInTheDocument();

      expect(getByText('Install plugin')).toBeInTheDocument();
      expect(getByText('kubectl krew install repo/packageName')).toBeInTheDocument();

      const link = getByText('Need Krew?');
      expect(link).toBeInTheDocument();
      expect(link).toHaveProperty('href', 'https://krew.sigs.k8s.io/docs/user-guide/setup/install/');
    });

    it('renders component when is default Krew repo', () => {
      const props = {
        ...defaultProps,
        repository: { ...defaultProps.repository, url: 'https://github.com/kubernetes-sigs/krew-index' },
      };
      const { getByText, queryByText } = render(<KrewInstall {...props} />);

      expect(queryByText('Add repository')).toBeNull();
      expect(queryByText('kubectl krew index add repo http://repo.test')).toBeNull();

      expect(getByText('Install plugin')).toBeInTheDocument();
      expect(getByText('kubectl krew install packageName')).toBeInTheDocument();

      const link = getByText('Need Krew?');
      expect(link).toBeInTheDocument();
      expect(link).toHaveProperty('href', 'https://krew.sigs.k8s.io/docs/user-guide/setup/install/');
    });
  });
});
