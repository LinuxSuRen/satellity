import React, { Component } from 'react';
import Config from '../components/config.js';
import style from './widget.scss';
import API from '../api/index.js';

class SiteWidget extends Component {
  constructor(props) {
    super(props);

    this.api = new API();
  }

  render() {
    let signIn;
    if (!this.api.user.loggedIn()) {
      signIn = (
        <div className={style.sign_in}>
          <a href={`https://github.com/login/oauth/authorize?scope=user:email&client_id=${Config.GithubClientId}`}>{i18n.t('login.github')}</a>
        </div>
      )
    }

    return (
      <div className={style.widget}>
        <div className={style.section}>
          <h2 className={style.site}>
            Go Discourse
          </h2>
          <ul className={style.features} dangerouslySetInnerHTML={{__html: i18n.t('aside.rules')}}>
          </ul>
        </div>
        {signIn}
        <div className={style.copyright}>
          © 2019 MIT license
        </div>
      </div>
    )
  }
}

export default SiteWidget;
