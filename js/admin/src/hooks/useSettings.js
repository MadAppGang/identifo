import { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchServerSetings } from '~/modules/settings/actions';
import { showSettingsSnack, hideSettingsSnack } from '~/modules/applications/actions';


export const useSettings = () => {
  const dispatch = useDispatch();
  const [changed, setChanged] = useState(false);
  const state = useSelector(s => s.settings);

  useEffect(() => {
    const fetch = async () => {
      await dispatch(fetchServerSetings());
    };
    fetch();
  }, []);

  useEffect(() => {
    const [original, current] = [JSON.stringify(state.original), JSON.stringify(state.current)];

    if (original !== current) {
      setChanged(true);
      dispatch(showSettingsSnack());
    }

    if (changed && original === current) {
      setChanged(false);
      dispatch(hideSettingsSnack());
    }
  }, [state, changed]);

  return { changed };
};
