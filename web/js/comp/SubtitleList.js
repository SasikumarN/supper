import React, { Component } from 'react';

import subtitleStore from '../stores/subtitle_store'

import FileTags from './FileTags'
import ScoreTag from './ScoreTag'

class SubtitleList extends Component {
  constructor() {
    super()
    this.state = subtitleStore.getState()
    this.updateSubtitles = this.updateSubtitles.bind(this)
  }

  updateSubtitles() {
    console.log("Update subtitles", subtitleStore.getSubtitles())
    this.setState(subtitleStore.getState())
  }

  componentWillMount() {
    subtitleStore.on("change", this.updateSubtitles)
  }

  componentWillUnmount() {
    subtitleStore.removeListener("change", this.updateSubtitles)
  }

  downloadSubtitle(sub) {
    console.log("download sub", sub)
    subtitleStore.download(sub)
  }

  render() {
    if (!this.state.subtitles || this.state.subtitles.length === 0) {
      return (
        <h3 className="meta center">Select a file to display subtitles</h3>
      )
    }

    let subs = this.state.subtitles

    subs = subs.filter((s) => s.language === this.state.lang)

    subs = subs.map((s) => {
      return (
        <li key={s.link} className="flex center collapse">
          <div className="col inline spaced flex center nowrap">
            <div className="col">
              <ScoreTag score={s.score}/>
            </div>
            <div className="col name ellipsis">{s.media.name}</div>
          </div>
          <div className="">
            <FileTags media={s}/>
          </div>
          <div className="right">
            <button className="small"
              onClick={() => this.downloadSubtitle(s)}>
              Download
            </button>
          </div>
        </li>
      )
    })

    return (
      <ul className="subtitle-list">{subs}</ul>
    )
  }
}

export default SubtitleList
