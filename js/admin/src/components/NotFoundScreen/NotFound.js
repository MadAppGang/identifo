import React from 'react';
import Header from '../ManagementScreen/Header';
import Container from '~/components/shared/Container';
import NotFound from '~/components/shared/NotFound';

const NotFoundScreen = () => (
  <div>
    <Header />
    <div className="iap-management-content">
      <Container>
        <div className="iap-global-not-found-screen">
          <NotFound />
        </div>
      </Container>
    </div>
  </div>
);

export default NotFoundScreen;
