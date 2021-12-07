import React, { CSSProperties } from 'react';

import { CanvasElementItem, CanvasElementProps } from '../element';
import {
  ColorDimensionConfig,
  ResourceDimensionConfig,
  ResourceDimensionMode,
  getPublicOrAbsoluteUrl,
} from 'app/features/dimensions';
import { ColorDimensionEditor, ResourceDimensionEditor } from 'app/features/dimensions/editors';
import SVG from 'react-inlinesvg';
import { css } from '@emotion/css';
import { isString } from 'lodash';
import { LineConfig } from '../types';
import { DimensionContext } from 'app/features/dimensions/context';
import { getBackendSrv } from '@grafana/runtime';
import { APIEditor, APIEditorConfig } from 'app/plugins/panel/canvas/editor/apiEditor';

export interface IconConfig {
  path?: ResourceDimensionConfig;
  fill?: ColorDimensionConfig;
  stroke?: LineConfig;
  api?: APIEditorConfig;
}

interface IconData {
  path: string;
  fill: string;
  strokeColor?: string;
  stroke?: number;
  api?: APIEditorConfig;
}

// When a stoke is defined, we want the path to be in page units
const svgStrokePathClass = css`
  path {
    vector-effect: non-scaling-stroke;
  }
`;

export function IconDisplay(props: CanvasElementProps) {
  const { width, height, data } = props;
  if (!data?.path) {
    return null;
  }

  const svgStyle: CSSProperties = {
    fill: data?.fill,
    stroke: data?.strokeColor,
    strokeWidth: data?.stroke,
  };

  const onClick = () => {
    if (data?.api) {
      getBackendSrv()
        .fetch({
          url: data?.api.endpoint!,
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Methods': 'POST',
            'Access-Control-Allow-Headers': 'Content-Type, Authorization',
          },
          data: data?.api.data ?? {},
        })
        .subscribe({
          next: (v: any) => {
            console.log('GOT', v);
          },
          error: (err: any) => {
            console.log('GOT ERROR', err);
            alert('TODO... button click: ' + JSON.stringify(err));
          },
          complete: () => {
            // this.setState({ working: false });
          },
        });
    }
  };

  return (
    <SVG
      onClick={onClick}
      src={data.path}
      width={width}
      height={height}
      style={svgStyle}
      className={svgStyle.strokeWidth ? svgStrokePathClass : undefined}
    />
  );
}

export const iconItem: CanvasElementItem<IconConfig, IconData> = {
  id: 'icon',
  name: 'Icon',
  description: 'SVG Icon display',

  display: IconDisplay,

  getNewOptions: (options) => ({
    placement: {
      width: 50,
      height: 50,
    },
    ...options,
    config: {
      path: {
        mode: ResourceDimensionMode.Fixed,
        fixed: 'img/icons/unicons/question-circle.svg',
      },
      fill: { fixed: '#FFF899' },
    },
  }),

  // Called when data changes
  prepareData: (ctx: DimensionContext, cfg: IconConfig) => {
    let path: string | undefined = undefined;
    if (cfg.path) {
      path = ctx.getResource(cfg.path).value();
    }
    if (!path || !isString(path)) {
      path = getPublicOrAbsoluteUrl('img/icons/unicons/question-circle.svg');
    }

    const data: IconData = {
      path,
      fill: cfg.fill ? ctx.getColor(cfg.fill).value() : '#CCC',
      api: cfg?.api ?? undefined,
    };

    if (cfg.stroke?.width && cfg.stroke.color) {
      if (cfg.stroke.width > 0) {
        data.stroke = cfg.stroke?.width;
        data.strokeColor = ctx.getColor(cfg.stroke.color).value();
      }
    }
    return data;
  },

  // Heatmap overlay options
  registerOptionsUI: (builder) => {
    const category = ['Icon'];
    builder
      .addCustomEditor({
        category,
        id: 'iconSelector',
        path: 'config.path',
        name: 'SVG Path',
        editor: ResourceDimensionEditor,
        settings: {
          resourceType: 'icon',
        },
      })
      .addCustomEditor({
        category,
        id: 'config.fill',
        path: 'config.fill',
        name: 'Fill color',
        editor: ColorDimensionEditor,
        settings: {},
        defaultValue: {
          // Configured values
          fixed: 'grey',
        },
      })
      .addSliderInput({
        category,
        path: 'config.stroke.width',
        name: 'Stroke',
        defaultValue: 0,
        settings: {
          min: 0,
          max: 10,
        },
      })
      .addCustomEditor({
        category,
        id: 'config.stroke.color',
        path: 'config.stroke.color',
        name: 'Stroke color',
        editor: ColorDimensionEditor,
        settings: {},
        defaultValue: {
          // Configured values
          fixed: 'grey',
        },
        showIf: (cfg) => Boolean(cfg?.config?.stroke?.width),
      })
      .addCustomEditor({
        category,
        id: 'apiSelector',
        path: 'config.api',
        name: 'API',
        editor: APIEditor,
      });
  },
};
