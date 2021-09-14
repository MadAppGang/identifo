import { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { getVerificationStatus } from '~/modules/settings/selectors';
import { setVerificationStatus, verifyConnection } from '~/modules/settings/actions';
import { verificationStatuses } from '~/enums';


export const useVerification = () => {
  const dispatch = useDispatch();
  const verificationStatus = useSelector(getVerificationStatus);
  const setStatus = async (status) => {
    await dispatch(setVerificationStatus(status));
  };
  useEffect(() => {
    return () => {
      dispatch(setVerificationStatus(verificationStatuses.required));
    };
  }, []);
  return [verificationStatus, verifyConnection, setStatus];
};
